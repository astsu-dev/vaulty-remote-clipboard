package services

import (
	"time"

	"golang.design/x/clipboard"
)

type ClipboardService struct {
	cancelClipboardCleanupChan chan struct{}
}

func (cs *ClipboardService) SetClipboard(content string) {
	clipboard.Write(clipboard.FmtText, []byte(content))
}

// Schedules clipboard cleanup after the specified timeout in seconds.
func (cs *ClipboardService) ScheduleClearClipboard(timeout time.Duration) {
	cs.cancelClipboardCleanupChan <- struct{}{}
	go func() {
		select {
		case <-cs.cancelClipboardCleanupChan:
		case <-time.After(timeout * time.Second):
			cs.SetClipboard("")
		}
	}()
}

func NewClipboardService() *ClipboardService {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	return &ClipboardService{}
}
