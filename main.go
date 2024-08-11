package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"remclip/internal/servers/udp"
	"remclip/internal/services/clipboard"
	"remclip/internal/services/clipboardapi"
	"remclip/internal/utils"
)

const (
	AppId                  = "dev.astsu.vaulty-remote-clipboard"
	WindowName             = "Vaulty Remote Clipboard"
	WindowWidth            = 300
	WindowHeight           = 300
	ServerStoppedLabelText = "The server is stopped"
	DefaultPort            = 8090
)

type ServerState struct {
	stop func()
}

func startServer(
	a fyne.App,
	pref fyne.Preferences,
	passwordEntry *widget.Entry,
	serverStatusLabel *widget.Label,
	serverState *ServerState,
	clipboardService clipboard.ClipboardServiceInterface,
) {
	if serverState.stop != nil {
		return
	}

	password := strings.TrimSpace(passwordEntry.Text)
	passwordEntry.SetText("")
	if password == "" {
		a.SendNotification(fyne.NewNotification("Invalid password", "Password can't be empty"))
		return
	}
	encryptionKey := utils.DerivePbkdf2From([]byte(password))

	port := pref.Int("port")

	server := udp.NewUDPServer(port, encryptionKey, clipboardService)
	ctx, cancel := context.WithCancel(context.Background())
	serverState.stop = cancel

	go func() {
		serverStatusLabel.SetText(getServerRunningText(port))
		err := server.Start(ctx)
		serverState.stop = nil
		if err != nil {
			a.SendNotification(fyne.NewNotification("Unexpected error", err.Error()))
		}
		serverStatusLabel.SetText(ServerStoppedLabelText)
	}()
}

func getServerRunningText(port int) string {
	return fmt.Sprintf("The server is running on port %d", port)
}

func main() {
	clipboardAPI, err := clipboardapi.NewClipClipboardAPI()
	if err != nil {
		panic(err)
	}
	clipboardService := clipboard.NewClipboardService(clipboardAPI)

	serverState := &ServerState{}

	a := app.NewWithID(AppId)
	pref := a.Preferences()

	w := a.NewWindow(WindowName)
	w.Resize(fyne.NewSize(WindowWidth, WindowHeight))

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port")
	savedPort := pref.IntWithFallback("port", DefaultPort)
	pref.SetInt("port", savedPort)
	portEntry.SetText(strconv.Itoa(savedPort))

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")
	serverStatusLabel := widget.NewLabel(ServerStoppedLabelText)
	serverStatusLabelContainer := container.NewCenter(serverStatusLabel)
	saveButton := widget.NewButton("Save", func() {
		port, err := strconv.Atoi(portEntry.Text)
		if err == nil {
			pref.SetInt("port", port)
		}
	})
	startButton := widget.NewButton("Start", func() {
		startServer(a, pref, passwordEntry, serverStatusLabel, serverState, clipboardService)
	})
	stopButton := widget.NewButton("Stop", func() {
		if serverState.stop != nil {
			serverState.stop()
		}
	})
	buttonsContainer := container.NewGridWithColumns(2, stopButton, startButton)
	content := container.NewVBox(portEntry, passwordEntry, saveButton, buttonsContainer, serverStatusLabelContainer)

	// Setup system tray
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu(
			WindowName,
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		)
		desk.SetSystemTrayMenu(m)

		w.SetCloseIntercept(func() {
			w.Hide()
		})
	}

	w.SetContent(content)
	w.ShowAndRun()
}
