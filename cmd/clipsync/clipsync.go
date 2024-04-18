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
)

type Config struct {
	ApiKeyForGetClipboard          string `toml:"apiKeyForGetClipboard"`
	ApiKeyForSetClipboard          string `toml:"apiKeyForsetClipboard"`
	HttpsCertificateFilePath       string `toml:"httpsCertificateFilePath"`
	HttpsCertificateSecretFilePath string `toml:"httpsCertificateSecretFilePath"`
	Port                           int    `toml:"port"`
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

	clipboardService := services.NewClipboardService()
	controller := controllers.Controller{
		ApiKeyForGetClipboard: config.ApiKeyForGetClipboard,
		ApiKeyForSetClipboard: config.ApiKeyForSetClipboard,
		ClipboardService:      clipboardService,
	}

	http.HandleFunc("GET /clipboard", controller.GetClipboardToSync)
	http.HandleFunc("POST /clipboard", controller.SetClipboardToSync)

	fmt.Printf("Your local address for client: https://%s:%d\n", getLocalIpAddress(), config.Port)
	addr := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Listening on https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(
		addr,
		config.HttpsCertificateFilePath,
		config.HttpsCertificateSecretFilePath,
		nil,
	))
}
