package services

import (
	"time"

	"golang.design/x/clipboard"
)

type ClipboardService struct{}

func (cs *ClipboardService) SetClipboard(content string) {
	clipboard.Write(clipboard.FmtText, []byte(content))
}

// Schedules clipboard cleanup after the specified timeout in seconds.
func (cs *ClipboardService) ScheduleClearClipboard(timeout time.Duration) {
	go func() {
		time.Sleep(timeout * time.Second)
		cs.SetClipboard("")
	}()
}

func NewClipboardService() *ClipboardService {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	return &ClipboardService{}
}
