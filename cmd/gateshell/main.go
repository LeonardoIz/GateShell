package main

import (
	"fmt"
	"gateshell/internal/proxy"
	"gateshell/internal/utils"
	"log"
	"path/filepath"
	"strconv"
)

const (
	serverName    = "GateShell"
	serverVersion = "0.1.0"
)

func PrintBanner() {
	// Print server banner
	fmt.Printf(`
%s v%s - A modern reverse proxy for SSH.

`, serverName, serverVersion)
}

func main() {
	PrintBanner()

	// Initialize configuration manager
	configManager := utils.NewManager()
	if err := configManager.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := configManager.GetConfig()

	// Initialize logging
	if err := utils.InitLogging(filepath.Join("data", cfg.Server.LogDir)); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}

	server := proxy.NewServer(&proxy.ServerConfig{
		Port:          strconv.Itoa(cfg.Server.Port),
		HostKeyFile:   filepath.Join("data", cfg.Server.HostKey),
		ServerName:    serverName,
		ServerVersion: serverVersion,
		Config:        configManager,
		AuthMethod:    cfg.Server.AuthMethod,
	})

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
