package udp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"remclip/internal/servers"
	"remclip/internal/utils/crypto"

	"github.com/go-playground/validator/v10"
)

// The interface for mocking the clipboard service in tests
type ClipboardServiceInterface interface {
	SetClipboard(content string)
	ScheduleClearClipboard(ctx context.Context, timeout uint)
}

type UDPServer struct {
	port             int
	encryptionKey    []byte
	clipboardService ClipboardServiceInterface
}

func (s *UDPServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)
	packetConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer packetConn.Close()
	log.Printf("Listening on %s for UDP messages\n", addr)

	// Run goroutine to handle server stopping with context
	go func() {
		<-ctx.Done()
		packetConn.Close()
	}()

	for {
		buf := make([]byte, 1024)

		n, _, err := packetConn.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		s.handleSetClipboardMessage(ctx, buf[:n])
	}
}

func (s *UDPServer) handleSetClipboardMessage(ctx context.Context, message []byte) {
	decryptedData, err := s.decryptBody(message)
	if err != nil {
		log.Println(err)
		return
	}

	var data *servers.SetClipboardBody
	if err := json.Unmarshal(decryptedData, &data); err != nil {
		log.Println(err)
		return
	}
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		log.Println(err)
		return
	}

	if data.ExpiresIn != nil && *data.ExpiresIn <= 0 {
		log.Println("Expiration time must be greater than 0")
		return
	}

	s.clipboardService.SetClipboard(*data.Text)

	if data.ExpiresIn != nil {
		s.clipboardService.ScheduleClearClipboard(ctx, uint(*data.ExpiresIn))
	}
}

// Parses and decrypts request body.
func (s *UDPServer) decryptBody(message []byte) ([]byte, error) {
	var body *servers.RequestBody
	if err := json.Unmarshal(message, &body); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return nil, err
	}

	decryptedData, err := crypto.DecryptGCM(s.encryptionKey, []byte(*body.Data))
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func NewUDPServer(
	port int,
	encryptionKey []byte,
	clipboardService ClipboardServiceInterface,
) *UDPServer {
	s := UDPServer{
		port:             port,
		encryptionKey:    encryptionKey,
		clipboardService: clipboardService,
	}
	return &s
}
