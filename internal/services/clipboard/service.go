package clipboard

import (
	"context"
	"sync"
	"time"
)

// The interface with Write method for writing the text content to the system clipboard
type ClipboardAPI interface {
	Write(content string)
}

// ClipboardService has methods for writing the text content to the system clipboard
// and scheduling the clipboard clearing after the specified timeout
type ClipboardService struct {
	clipboardAPI ClipboardAPI

	// The mutex for the callback for cancelling previously scheduler clear clipboard goroutine
	m sync.Mutex
	// The function to cancel the previous scheduled clear clipboard goroutine
	// before a new one will be scheduled
	cancelScheduledClearClipboard *context.CancelFunc
}

func (cs *ClipboardService) SetClipboard(content string) {
	cs.clipboardAPI.Write(content)
}

// Schedules clipboard cleanup after the specified timeout in seconds.
// ctx can be used to cancel the scheduled goroutine.
func (cs *ClipboardService) ScheduleClearClipboard(
	ctx context.Context,
	timeout uint,
) {
	cs.m.Lock()
	if cs.cancelScheduledClearClipboard != nil {
		(*cs.cancelScheduledClearClipboard)()
	}

	ctx, cancel := context.WithCancel(ctx)
	cs.cancelScheduledClearClipboard = &cancel
	cs.m.Unlock()

	go func() {
		select {
		case <-time.After(time.Duration(timeout) * time.Second):
			cs.SetClipboard("")
			cs.m.Lock()
			cs.cancelScheduledClearClipboard = nil
			cs.m.Unlock()
		case <-ctx.Done():
			return
		}
	}()
}

func NewClipboardService(clipboardAPI ClipboardAPI) *ClipboardService {
	return &ClipboardService{
		clipboardAPI: clipboardAPI,
	}
}
