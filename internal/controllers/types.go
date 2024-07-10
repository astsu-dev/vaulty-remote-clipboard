package controllers

type RequestBody struct {
	Data string `json:"data"`
}

type SetClipboardBody struct {
	Text      string `json:"text"`
	ExpiresIn *int64 `json:"expiresIn"`
}
