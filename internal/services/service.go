package services

import (
	"golang.design/x/clipboard"
)

type ClipboardService struct {
	ClipboardToSync string
}

func (cs *ClipboardService) GetClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func NewClipboardService() *ClipboardService {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	return &ClipboardService{
		ClipboardToSync: "",
	}
}
