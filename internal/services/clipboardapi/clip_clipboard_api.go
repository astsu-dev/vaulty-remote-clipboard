package clipboardapi

import "golang.design/x/clipboard"

type ClipClipboardAPI struct{}

func (ca *ClipClipboardAPI) Write(content string) {
	clipboard.Write(clipboard.FmtText, []byte(content))
}

func NewClipClipboardAPI() (*ClipClipboardAPI, error) {
	err := clipboard.Init()
	if err != nil {
		return nil, err
	}
	return &ClipClipboardAPI{}, nil
}
