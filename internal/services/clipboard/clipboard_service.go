package clipboard

import (
	"context"
	"time"
)

// The interface with Write method for writing the text content to the system clipboard
type ClipboardAPI interface {
	Write(content string)
}

// The interface for mocking the clipboard service in tests
type ClipboardServiceInterface interface {
	SetClipboard(content string)
	ScheduleClearClipboard(ctx context.Context, timeout uint)
}

// ClipboardService has methods for writing the text content to the system clipboard
// and scheduling the clipboard clearing after the specified timeout
type ClipboardService struct {
	clipboardAPI ClipboardAPI

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
	if cs.cancelScheduledClearClipboard != nil {
		(*cs.cancelScheduledClearClipboard)()
	}

	ctx, cancel := context.WithCancel(ctx)
	cs.cancelScheduledClearClipboard = &cancel

	go func() {
		select {
		case <-time.After(time.Duration(timeout) * time.Second):
			cs.SetClipboard("")
			cs.cancelScheduledClearClipboard = nil
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
