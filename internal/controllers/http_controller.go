package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"clipsync/internal/services"
	"clipsync/internal/utils"
)

type HTTPController struct {
	EncryptionKey    []byte
	ClipboardService *services.ClipboardService
}

func (c *HTTPController) SetClipboard(w http.ResponseWriter, req *http.Request) {
	decryptedData, skip := c.decryptBody(w, req)
	if skip {
		return
	}

	var data *SetClipboardBody
	if err := json.Unmarshal(decryptedData, &data); err != nil {
		c.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Bad Request"})
		return
	}

	if data.ExpiresIn != nil && *data.ExpiresIn <= 0 {
		c.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Expiration time must be greater than 0"})
		return
	}

	c.ClipboardService.SetClipboard(data.Text)

	if data.ExpiresIn != nil {
		c.ClipboardService.ScheduleClearClipboard(time.Duration(*data.ExpiresIn))
	}

	c.sendResponse(w, http.StatusOK, map[string]any{})
}

func (c *HTTPController) sendResponse(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	serializedData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	encrypted := utils.EncryptGCM(c.EncryptionKey, serializedData)

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
func (c *HTTPController) decryptBody(w http.ResponseWriter, req *http.Request) (result []byte, skip bool) {
	var body *RequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		c.sendResponse(w, http.StatusBadRequest, map[string]any{"error": "Bad Request"})
		return nil, true
	}

	decryptedData := utils.DecryptGCM(c.EncryptionKey, []byte(body.Data))
	return decryptedData, false
}
