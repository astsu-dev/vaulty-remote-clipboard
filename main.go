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

	"remclip/internal/servers"
	"remclip/internal/utils"
)

const (
	APP_ID                    = "dev.astsu.vaulty-remote-clipboard"
	WINDOW_NAME               = "Vaulty Remote Clipboard"
	WINDOW_WIDTH              = 300
	WINDOW_HEIGHT             = 300
	SERVER_STOPPED_LABEL_TEXT = "The server is stopped"
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

	server := servers.NewUDPServer(port, encryptionKey)
	ctx, cancel := context.WithCancel(context.Background())
	serverState.stop = cancel

	go func() {
		serverStatusLabel.SetText(getServerRunningText(port))
		err := server.Start(ctx)
		serverState.stop = nil
		if err != nil {
			a.SendNotification(fyne.NewNotification("Unexpected error", err.Error()))
		}
		serverStatusLabel.SetText(SERVER_STOPPED_LABEL_TEXT)
	}()
}

func getServerRunningText(port int) string {
	return fmt.Sprintf("The server is running on port %d", port)
}

func main() {
	serverState := &ServerState{}

	a := app.NewWithID(APP_ID)
	pref := a.Preferences()

	w := a.NewWindow(WINDOW_NAME)
	w.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT))

	portEntry := widget.NewEntry()
	savedPort := pref.IntWithFallback("port", 8090)
	pref.SetInt("port", savedPort)
	portEntry.SetText(strconv.Itoa(savedPort))

	passwordEntry := widget.NewPasswordEntry()
	serverStatusLabel := widget.NewLabel("The server is stopped")
	serverStatusLabelContainer := container.NewCenter(serverStatusLabel)
	saveButton := widget.NewButton("Save", func() {
		port, err := strconv.Atoi(portEntry.Text)
		if err == nil {
			pref.SetInt("port", port)
		}
	})
	startButton := widget.NewButton("Start", func() {
		startServer(a, pref, passwordEntry, serverStatusLabel, serverState)
	})
	stopButton := widget.NewButton("Stop", func() {
		if serverState.stop != nil {
			serverState.stop()
		}
	})
	buttonsContainer := container.NewGridWithColumns(2, startButton, stopButton)
	content := container.NewVBox(portEntry, passwordEntry, saveButton, buttonsContainer, serverStatusLabelContainer)

	// Setup system tray
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu(
			WINDOW_NAME,
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
