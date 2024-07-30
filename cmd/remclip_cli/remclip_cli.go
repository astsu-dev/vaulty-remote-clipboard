package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/term"

	"remclip/internal/servers"
	"remclip/internal/utils"
)

const (
	HTTPServerType = "http"
	UDPServerType  = "udp"
)

type Config struct {
	Port       int    `toml:"port"`
	ServerType string `toml:"serverType"`
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
		panic("Port must be greater than 0")
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
		configPath = path.Join(homeDir, ".config", "remclip", "config.toml")
	default:
		configPath = "config.toml"
	}

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	return configPath
}

func main() {
	config := loadConfig(resolveConfigPath())

	fmt.Printf("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Println()

	encryptionKey := utils.DerivePbkdf2From(password)

	var server servers.Server
	switch config.ServerType {
	case HTTPServerType:
		server = servers.NewHTTPServer(config.Port, encryptionKey)
	case UDPServerType:
		server = servers.NewUDPServer(config.Port, encryptionKey)
	default:
		log.Fatalf("Unexpected server type: %s", config.ServerType)
	}

	// TODO: Use cancel context instead when server stopping will be supported
	ctx := context.TODO()

	err = server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
