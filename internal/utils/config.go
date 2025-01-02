package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultConfigFile = "config.json"
	DefaultPort       = 22
	DefaultHostKey    = "ssh_host_key"
	DefaultLogDir     = "logs"
)

// AuthConfig represents the authentication configuration for an endpoint
type AuthConfig struct {
	User    string   `json:"user"`
	Methods []string `json:"methods"`
}

// EndpointConfig represents the configuration of an endpoint
type EndpointConfig struct {
	Target string     `json:"target"`
	Auth   AuthConfig `json:"auth"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port            int    `json:"port"`
	HostKey         string `json:"host_key"`
	DefaultEndpoint string `json:"default_endpoint,omitempty"`
	LogDir          string `json:"log_dir"`
}

// Config represents the complete configuration
type Config struct {
	Server    ServerConfig              `json:"server"`
	Endpoints map[string]EndpointConfig `json:"endpoints,omitempty"`
}

// Manager manages the configuration
type Manager struct {
	configFile string
	config     *Config
}

// NewManager creates a new configuration manager instance with the given config file path
func NewManager(configFile string) *Manager {
	return &Manager{
		configFile: configFile,
	}
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    DefaultPort,
			HostKey: DefaultHostKey,
			LogDir:  DefaultLogDir,
		},
	}
}

// LoadConfig loads the configuration from the file
func (m *Manager) LoadConfig() error {
	// Check if the file exists
	if _, err := os.Stat(m.configFile); os.IsNotExist(err) {
		// Create default configuration
		m.config = createDefaultConfig()
		// Save default configuration
		if err := m.SaveConfig(); err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
		return nil
	}

	// Read configuration file
	data, err := os.ReadFile(m.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Decode JSON
	m.config = &Config{}
	if err := json.Unmarshal(data, m.config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate configuration
	if err := m.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

// SaveConfig saves the configuration to the file
func (m *Manager) SaveConfig() error {
	// Create directory if it does not exist
	dir := filepath.Dir(m.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Encode JSON
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %v", err)
	}

	// Write file
	if err := os.WriteFile(m.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// validateConfig validates the configuration
func (m *Manager) validateConfig() error {
	// Validate port
	if m.config.Server.Port <= 0 || m.config.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", m.config.Server.Port)
	}

	// Validate host key
	if m.config.Server.HostKey == "" {
		return fmt.Errorf("host key file path is required")
	}

	// Validate log directory
	if m.config.Server.LogDir == "" {
		return fmt.Errorf("log directory path is required")
	}

	// Validate endpoints
	if m.config.Endpoints != nil {
		// Validate default endpoint
		if m.config.Server.DefaultEndpoint != "" {
			if _, exists := m.config.Endpoints[m.config.Server.DefaultEndpoint]; !exists {
				return fmt.Errorf("default endpoint '%s' not found", m.config.Server.DefaultEndpoint)
			}
		}

		// Validate each endpoint
		validMethods := map[string]bool{
			"password": true,
			"none":     true,
		}

		for name, endpoint := range m.config.Endpoints {
			if endpoint.Target == "" {
				return fmt.Errorf("target is required for endpoint '%s'", name)
			}

			if endpoint.Auth.User == "" {
				return fmt.Errorf("user is required for endpoint '%s'", name)
			}

			for _, method := range endpoint.Auth.Methods {
				if !validMethods[method] {
					return fmt.Errorf("invalid auth method '%s' for endpoint '%s'", method, name)
				}
			}
		}
	}

	return nil
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *Config {
	return m.config
}

// GetEndpoint gets the configuration of a specific endpoint
func (m *Manager) GetEndpoint(name string) (*EndpointConfig, error) {
	if m.config.Endpoints == nil {
		return nil, fmt.Errorf("no endpoints configured")
	}

	endpoint, exists := m.config.Endpoints[name]
	if !exists {
		if m.config.Server.DefaultEndpoint == "" {
			return nil, fmt.Errorf("endpoint '%s' not found and no default endpoint configured", name)
		}
		endpoint = m.config.Endpoints[m.config.Server.DefaultEndpoint]
	}

	return &endpoint, nil
}
