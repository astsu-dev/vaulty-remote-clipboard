package servers

type RequestBody struct {
	Data *string `json:"data" validate:"required"`
}

type SetClipboardBody struct {
	Text      *string `json:"text" validate:"required"`
	ExpiresIn *int64  `json:"expiresIn" validate:"omitempty,gte=0"`
}
