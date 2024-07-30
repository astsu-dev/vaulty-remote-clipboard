package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"remclip/internal/services"
	"remclip/internal/utils"
	"time"
)

type HTTPServer struct {
	encryptionKey    []byte
	server           *http.Server
	clipboardService *services.ClipboardService
}

func (s *HTTPServer) Start(ctx context.Context) error {
	listenResultCh := make(chan error)

	go func() {
		err := s.server.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
		listenResultCh <- err
	}()

	select {
	case <-ctx.Done():
		err := s.server.Shutdown(context.Background())
		if err != nil {
			return err
		}
		err = <-listenResultCh
		if err != nil {
			return err
		}
	case err := <-listenResultCh:
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HTTPServer) setClipboardHandler(w http.ResponseWriter, req *http.Request) {
	decryptedData, skip := s.decryptBody(w, req)
	if skip {
		return
	}

	var data *SetClipboardBody
	if err := json.Unmarshal(decryptedData, &data); err != nil {
		s.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Bad Request"})
		return
	}

	if data.ExpiresIn != nil && *data.ExpiresIn <= 0 {
		s.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Expiration time must be greater than 0"})
		return
	}

	s.clipboardService.SetClipboard(data.Text)

	if data.ExpiresIn != nil {
		s.clipboardService.ScheduleClearClipboard(time.Duration(*data.ExpiresIn))
	}

	s.sendResponse(w, http.StatusOK, map[string]any{})
}

func (s *HTTPServer) sendResponse(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	serializedData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	encrypted := utils.EncryptGCM(s.encryptionKey, serializedData)

	serializedBody, err := json.Marshal(map[string]any{
		"data": string(encrypted),
	})
	if err != nil {
		panic(err)
	}

	_, err = w.Write(serializedBody)
	if err != nil {
		panic(err)
	}
}

// Parses and decrypts request body.
func (s *HTTPServer) decryptBody(w http.ResponseWriter, req *http.Request) (result []byte, skip bool) {
	var body *RequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		s.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Bad Request"})
		return nil, true
	}

	decryptedData, err := utils.DecryptGCM(s.encryptionKey, []byte(body.Data))
	if err != nil {
		s.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Message authentication fails"})
		return nil, true
	}
	return decryptedData, false
}

func NewHTTPServer(port int, encryptionKey []byte) *HTTPServer {
	clipboardService := services.NewClipboardService()
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()
	httpServer := http.Server{Addr: addr, Handler: mux}
	s := HTTPServer{
		encryptionKey:    encryptionKey,
		server:           &httpServer,
		clipboardService: clipboardService,
	}
	http.HandleFunc("POST /clipboard", s.setClipboardHandler)

	return &s
}
