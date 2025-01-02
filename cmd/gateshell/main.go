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
	// Print the server banner with name and version
	fmt.Printf(`
%s v%s - A modern reverse proxy for SSH.

`, serverName, serverVersion)
}

func main() {
	PrintBanner()

	// Initialize configuration manager
	configManager := config.NewManager()
	// Load configuration from file
	if err := configManager.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Retrieve server configuration
	cfg := configManager.GetConfig()
	// Create a new proxy server with the loaded configuration
	server := proxy.NewServer(&proxy.ServerConfig{
		Port:          strconv.Itoa(cfg.Server.Port),
		HostKeyFile:   cfg.Server.HostKey,
		ServerName:    serverName,
		ServerVersion: serverVersion,
		Config:        configManager,
	})

	// Start the proxy server
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
