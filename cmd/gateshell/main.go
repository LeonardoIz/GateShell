// cmd/gateshell/main.go
package main

import (
	"fmt"
	"gateshell/internal/config"
	"gateshell/internal/proxy"
	"log"
	"strconv"
)

const (
	serverName    = "GateShell"
	serverVersion = "0.1.0"
)

func PrintBanner() {
	fmt.Printf(`
%s v%s - A modern reverse proxy for SSH.

`, serverName, serverVersion)
}

func main() {
	PrintBanner()

	// Cargar configuración
	configManager := config.NewManager()
	if err := configManager.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := configManager.GetConfig()
	server := proxy.NewServer(&proxy.ServerConfig{
		Port:          strconv.Itoa(cfg.Server.Port),
		HostKeyFile:   cfg.Server.HostKey,
		ServerName:    serverName,
		ServerVersion: serverVersion,
		Config:        configManager,
	})

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
