package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"clipsync/internal/services"
	"clipsync/internal/utils"
)

type UDPController struct {
	EncryptionKey    []byte
	ClipboardService *services.ClipboardService
}

func (c *UDPController) SetClipboard(message []byte) {
	decryptedData, err := c.decryptBody(message)
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

	c.ClipboardService.SetClipboard(data.Text)

	if data.ExpiresIn != nil {
		c.ClipboardService.ScheduleClearClipboard(time.Duration(*data.ExpiresIn))
	}
}

// Parses and decrypts request body.
func (c *UDPController) decryptBody(message []byte) ([]byte, error) {
	var body *RequestBody
	if err := json.Unmarshal(message, &body); err != nil {
		return nil, errors.New("Invalid UDP message")
	}

	decryptedData := utils.DecryptGCM(c.EncryptionKey, []byte(body.Data))
	return decryptedData, nil
}
