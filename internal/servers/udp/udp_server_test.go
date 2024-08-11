package udp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"remclip/internal/utils"
	"testing"
	"time"
)

type fakeClipboardService struct {
	setClipboardCalls           []string
	scheduleClearClipboardCalls []uint
}

func (s *fakeClipboardService) SetClipboard(content string) {
	s.setClipboardCalls = append(s.setClipboardCalls, content)
}

func (s *fakeClipboardService) ScheduleClearClipboard(ctx context.Context, timeout uint) {
	s.scheduleClearClipboardCalls = append(s.scheduleClearClipboardCalls, timeout)
}

func (s *fakeClipboardService) SetClipboardHasBeenCalledWith(content string) bool {
	for _, c := range s.setClipboardCalls {
		if c == content {
			return true
		}
	}
	return false
}

func (s *fakeClipboardService) SetClipboardHasBeenCalled() bool {
	return len(s.setClipboardCalls) > 0
}

func (s *fakeClipboardService) ScheduleClearClipboardHasBeenCalledWith(timeout uint) bool {
	for _, t := range s.scheduleClearClipboardCalls {
		if t == timeout {
			return true
		}
	}
	return false
}

func (s *fakeClipboardService) ScheduleClearClipboardHasBeenCalled() bool {
	return len(s.scheduleClearClipboardCalls) > 0
}

func TestUDPServer(t *testing.T) {
	port := 9123
	connectionString := fmt.Sprintf(":%d", port)
	encryptionKey := utils.DerivePbkdf2From([]byte("testpassword"))

	t.Run("should set clipboard and schedule clear after the request", func(t *testing.T) {
		// given
		expectedClipboardContent := "test"
		var expectedClearTimeout uint = 1
		clipboardService := &fakeClipboardService{}
		server := NewUDPServer(port, encryptionKey, clipboardService)

		// stop the server at the end of the test
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// run the server
		go func() {
			err := server.Start(ctx)
			if err != nil {
				t.Logf("server is stopped with the error: %v", err)
			}
		}()

		// when
		// make request
		con, err := net.Dial("udp", connectionString)
		if err != nil {
			t.Fatalf("cannot connect to the server %s: %v", connectionString, err)
		}
		defer con.Close()
		sendRequest(
			t, con, encryptionKey,
			map[string]any{"text": expectedClipboardContent, "expiresIn": expectedClearTimeout},
		)
		// wait until the request will be processed
		time.Sleep(100 * time.Millisecond)

		// then
		if !clipboardService.SetClipboardHasBeenCalledWith(expectedClipboardContent) {
			t.Fatalf("SetClipboard was not called with the %s", expectedClipboardContent)
		}
		if !clipboardService.ScheduleClearClipboardHasBeenCalledWith(expectedClearTimeout) {
			t.Fatalf("ScheduleSetClipboard was not called with timeout %d", expectedClearTimeout)
		}
	})

	t.Run("should not schedule clear if expiresIn not passed in the request body", func(t *testing.T) {
		// given
		expectedClipboardContent := "test"
		clipboardService := &fakeClipboardService{}
		server := NewUDPServer(port, encryptionKey, clipboardService)

		// stop the server at the end of the test
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// run the server
		go func() {
			err := server.Start(ctx)
			if err != nil {
				t.Logf("server is stopped with the error: %v", err)
			}
		}()

		// when
		// make request
		con, err := net.Dial("udp", connectionString)
		if err != nil {
			t.Fatalf("cannot connect to the server %s: %v", connectionString, err)
		}
		defer con.Close()
		sendRequest(
			t, con, encryptionKey,
			map[string]any{"text": expectedClipboardContent},
		)
		// wait until the request will be processed
		time.Sleep(100 * time.Millisecond)

		// then
		if clipboardService.ScheduleClearClipboardHasBeenCalled() {
			t.Fatal("ScheduleSetClipboard was called, but must not")
		}
		if !clipboardService.SetClipboardHasBeenCalledWith(expectedClipboardContent) {
			t.Fatalf("SetClipboard was not called with the %s", expectedClipboardContent)
		}
	})

	t.Run("should not send clipboard and not schedule clear if sent the invalid request body", func(t *testing.T) {
		testCases := []struct {
			Name string
			Body map[string]any
		}{
			{Name: "empty body", Body: map[string]any{}},
			{Name: "negative expiresIn", Body: map[string]any{"text": "test", "expiresIn": -1}},
		}
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				// given
				clipboardService := &fakeClipboardService{}
				server := NewUDPServer(port, encryptionKey, clipboardService)

				// stop the server at the end of the test
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				// run the server
				go func() {
					err := server.Start(ctx)
					if err != nil {
						t.Logf("server is stopped with the error: %v", err)
					}
				}()

				// when
				// make request
				con, err := net.Dial("udp", connectionString)
				if err != nil {
					t.Fatalf("cannot connect to the server %s: %v", connectionString, err)
				}
				defer con.Close()
				sendRequest(t, con, encryptionKey, map[string]any{})
				// wait until the request will be processed
				time.Sleep(100 * time.Millisecond)

				// then
				if clipboardService.SetClipboardHasBeenCalled() {
					t.Fatal("SetClipboard was called, but must not")
				}
				if clipboardService.ScheduleClearClipboardHasBeenCalled() {
					t.Fatal("ScheduleSetClipboard was called, but must not")
				}
			})
		}
	})
}

func sendRequest(
	t *testing.T, con net.Conn, encryptionKey []byte, data map[string]any,
) {
	encodedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("cannot marshal the set clipboard body: %v", err)
	}
	encryptedSetClipboardBody := utils.EncryptGCM(encryptionKey, encodedData)
	requestBody := map[string]string{
		"data": string(encryptedSetClipboardBody),
	}
	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("cannot marshal the request body: %v", err)
	}
	_, err = con.Write(encodedRequestBody)
	if err != nil {
		t.Fatalf("cannot send request to the server")
	}
}
