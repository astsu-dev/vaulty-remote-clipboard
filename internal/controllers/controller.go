package controllers

import (
	"encoding/json"
	"net/http"

	"clipsync/internal/services"
)

type Controller struct {
	ApiKeyForGetClipboard string
	ApiKeyForSetClipboard string
	ClipboardService      *services.ClipboardService
}

func (c *Controller) GetClipboardToSync(w http.ResponseWriter, req *http.Request) {
	if reqApiKeys, ok := req.Header["X-Api-Key"]; !ok || len(reqApiKeys) != 1 || reqApiKeys[0] != c.ApiKeyForGetClipboard {
		sendResponse(w, http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
		return
	}

	clipboard := c.ClipboardService.ClipboardToSync
	// Clear clipboard after the consumption by client
	c.ClipboardService.ClipboardToSync = ""

	sendResponse(w, http.StatusOK, map[string]any{
		"text": clipboard,
	})
}

func (c *Controller) SetClipboardToSync(w http.ResponseWriter, req *http.Request) {
	if reqApiKeys, ok := req.Header["X-Api-Key"]; !ok || len(reqApiKeys) != 1 || reqApiKeys[0] != c.ApiKeyForSetClipboard {
		sendResponse(w, http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
		return
	}

	c.ClipboardService.ClipboardToSync = c.ClipboardService.GetClipboard()

	sendResponse(w, http.StatusOK, map[string]any{})
}

func sendResponse(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	serialized, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(serialized)
	if err != nil {
		panic(err)
	}
}
