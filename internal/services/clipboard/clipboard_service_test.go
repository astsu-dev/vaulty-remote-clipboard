package clipboard

import (
	"context"
	"testing"
	"time"
)

type fakeClipboardAPI struct {
	writeCalls []string
}

func (ca *fakeClipboardAPI) Write(content string) {
	ca.writeCalls = append(ca.writeCalls, content)
}

func (ca *fakeClipboardAPI) WriteHasBeenCalled() bool {
	return len(ca.writeCalls) > 0
}

func (ca *fakeClipboardAPI) WriteHasBeenCalledTimes(n uint) bool {
	return len(ca.writeCalls) == int(n)
}

func (ca *fakeClipboardAPI) WriteHasBeenCalledWith(content string) bool {
	for _, c := range ca.writeCalls {
		if c == content {
			return true
		}
	}
	return false
}

func (ca *fakeClipboardAPI) WriteHasBeenCalledTimesWith(content string, n uint) bool {
	var count uint
	for _, c := range ca.writeCalls {
		if c == content {
			count++
		}
	}
	return count == n
}

func (ca *fakeClipboardAPI) WriteCalledTimes() uint {
	return uint(len(ca.writeCalls))
}

func TestClipboardService(t *testing.T) {
	t.Run("SetClipboard", func(t *testing.T) {
		t.Run("should call Write method of the given ClipboardAPI", func(t *testing.T) {
			contents := []string{"", "test"}

			for _, content := range contents {
				// given
				clipboardAPI := &fakeClipboardAPI{}
				service := NewClipboardService(clipboardAPI)

				// when
				service.SetClipboard(content)

				// then
				if !clipboardAPI.WriteHasBeenCalledWith(content) {
					t.Fatalf("Write method of the ClipboardAPI was not called with %s", content)
				}
			}
		})

	})

	t.Run("ScheduleClearClipboard", func(t *testing.T) {
		t.Run("should call Write method of the ClipboardAPI with empty string after the given timeout", func(t *testing.T) {
			timeouts := []uint{0, 1}

			for _, timeout := range timeouts {
				// given
				clipboardAPI := &fakeClipboardAPI{}
				service := NewClipboardService(clipboardAPI)
				expectedContent := ""
				var expectedCallTimes uint = 1

				// when
				service.ScheduleClearClipboard(context.Background(), timeout)
				if clipboardAPI.WriteHasBeenCalled() {
					t.Fatalf("Write method of the ClipboardAPI was called before the timeout")
				}
				time.Sleep(time.Duration(timeout+1) * time.Second)

				// then
				if !clipboardAPI.WriteHasBeenCalledTimes(expectedCallTimes) {
					t.Fatalf(
						"Write method of the ClipboardAPI was not called %d times, but %d instead",
						expectedCallTimes,
						clipboardAPI.WriteCalledTimes(),
					)
				}
				if !clipboardAPI.WriteHasBeenCalledWith(expectedContent) {
					t.Fatalf("Write method of the ClipboardAPI was not called with %s", expectedContent)
				}
			}
		})

		t.Run("should cancel the first scheduled clear clipboard after the second call and schedule the second", func(t *testing.T) {
			// given
			clipboardAPI := &fakeClipboardAPI{}
			service := NewClipboardService(clipboardAPI)
			emptyContent := ""
			var timeout1 uint = 1
			var timeout2 uint = 2

			// when
			// first call
			service.ScheduleClearClipboard(context.Background(), timeout1)
			if clipboardAPI.WriteHasBeenCalled() {
				t.Fatal("Write method of the ClipboardAPI was called before the timeout")
			}

			// second call
			service.ScheduleClearClipboard(context.Background(), timeout2)
			if clipboardAPI.WriteHasBeenCalled() {
				t.Fatal("Write method of the ClipboardAPI was called before the timeout")
			}

			// then
			// check the Write method was not called after the scheduled timeout as it was cancelled
			time.Sleep(time.Duration(timeout1) * time.Second)
			if clipboardAPI.WriteHasBeenCalled() {
				t.Fatal("Write method of the ClipboardAPI was called, but must not")
			}

			// wait the second scheduled goroutine to complete
			time.Sleep(time.Duration(timeout2) * time.Second)
			if !clipboardAPI.WriteHasBeenCalledTimes(1) {
				t.Fatalf(
					"Write method of the ClipboardAPI was called more than once, but %d instead",
					clipboardAPI.WriteCalledTimes(),
				)
			}
			if !clipboardAPI.WriteHasBeenCalledWith(emptyContent) {
				t.Fatalf("Write method of the ClipboardAPI was not called with %s", emptyContent)
			}
		})
	})
}
