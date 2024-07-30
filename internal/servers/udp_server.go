package servers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"remclip/internal/services"
	"remclip/internal/utils"
)

type UDPServer struct {
	port             int
	encryptionKey    []byte
	clipboardService *services.ClipboardService
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

		s.handleSetClipboardMessage(buf[:n])
	}
}

func (s *UDPServer) handleSetClipboardMessage(message []byte) {
	decryptedData, err := s.decryptBody(message)
	if err != nil {
		log.Println(err)
		return
	}

	var data *SetClipboardBody
	if err := json.Unmarshal(decryptedData, &data); err != nil {
		log.Println(err)
		return
	}

	if data.ExpiresIn != nil && *data.ExpiresIn <= 0 {
		log.Println("Expiration time must be greater than 0")
		return
	}

	s.clipboardService.SetClipboard(data.Text)

	if data.ExpiresIn != nil {
		s.clipboardService.ScheduleClearClipboard(time.Duration(*data.ExpiresIn))
	}
}

// Parses and decrypts request body.
func (s *UDPServer) decryptBody(message []byte) ([]byte, error) {
	var body *RequestBody
	if err := json.Unmarshal(message, &body); err != nil {
		return nil, errors.New("Invalid UDP message")
	}

	decryptedData, err := utils.DecryptGCM(s.encryptionKey, []byte(body.Data))
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func NewUDPServer(port int, encryptionKey []byte) *UDPServer {
	clipboardService := services.NewClipboardService()

	s := UDPServer{
		port:             port,
		encryptionKey:    encryptionKey,
		clipboardService: clipboardService,
	}

	return &s
}
