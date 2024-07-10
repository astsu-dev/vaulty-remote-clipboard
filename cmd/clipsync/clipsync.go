package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/pelletier/go-toml/v2"

	"clipsync/internal/controllers"
	"clipsync/internal/services"
	"clipsync/internal/utils"
)

const (
	HTTPServerType = "http"
	UDPServerType  = "udp"
)

type Config struct {
	EncryptionKey string `toml:"encryptionKey"`
	Port          int    `toml:"port"`
	ServerType    string `toml:"serverType"`
}

func loadConfig(path string) *Config {
	fmt.Printf("Loading config from %s\n", path)

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	config := &Config{}

	err = toml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	if config.Port < 1 {
		panic("You must specify correct port in the config")
	}

	switch config.ServerType {
	case HTTPServerType, UDPServerType:
	default:
		panic("You must specify valid server type. The valid values are: http or udp")
	}

	return config
}

func resolveConfigPath() string {
	var configPath string

	switch runtime.GOOS {
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configPath = path.Join(homeDir, ".config", "clipsync", "config.toml")
	default:
		configPath = "config.toml"
	}

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	return configPath
}

func getLocalIpAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}

func main() {
	config := loadConfig(resolveConfigPath())
	encryptionKey := utils.DerivePbkdf2From([]byte(config.EncryptionKey))
	clipboardService := services.NewClipboardService()

	switch config.ServerType {
	case HTTPServerType:
		controller := controllers.HTTPController{
			EncryptionKey:    encryptionKey,
			ClipboardService: clipboardService,
		}

		http.HandleFunc("POST /clipboard", controller.SetClipboard)

		fmt.Printf("Your local address for client: http://%s:%d\n", getLocalIpAddress(), config.Port)
		addr := fmt.Sprintf(":%d", config.Port)
		fmt.Printf("Listening on http://%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, nil))

	case UDPServerType:
		controller := controllers.UDPController{
			EncryptionKey:    encryptionKey,
			ClipboardService: clipboardService,
		}

		addr := fmt.Sprintf(":%d", config.Port)
		packetConn, err := net.ListenPacket("udp", addr)
		if err != nil {
			panic(err)
		}
		defer packetConn.Close()
		fmt.Printf("Listening on %s for UDP messages\n", addr)

		for {
			buf := make([]byte, 1024)
			n, _, err := packetConn.ReadFrom(buf)
			if err != nil {
				panic(err)
			}

			controller.SetClipboard(buf[:n])
		}
	}
}
