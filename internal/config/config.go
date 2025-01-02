// internal/config/config.go
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultConfigFile = "config.json"
	DefaultPort       = 22
	DefaultHostKey    = "ssh_host_key"
)

// AuthConfig representa la configuración de autenticación para un endpoint
type AuthConfig struct {
	User    string   `json:"user"`
	Methods []string `json:"methods"`
}

// EndpointConfig representa la configuración de un endpoint
type EndpointConfig struct {
	Target string     `json:"target"`
	Auth   AuthConfig `json:"auth"`
}

// ServerConfig representa la configuración del servidor
type ServerConfig struct {
	Port            int    `json:"port"`
	HostKey         string `json:"host_key"`
	DefaultEndpoint string `json:"default_endpoint,omitempty"`
}

// Config representa la configuración completa
type Config struct {
	Server    ServerConfig              `json:"server"`
	Endpoints map[string]EndpointConfig `json:"endpoints,omitempty"`
}

// Manager gestiona la configuración
type Manager struct {
	configFile string
	config     *Config
}

// NewManager crea una nueva instancia del gestor de configuración
func NewManager() *Manager {
	return &Manager{
		configFile: getConfigPath(),
	}
}

// getConfigPath obtiene la ruta del archivo de configuración
func getConfigPath() string {
	// Comprobar flag de línea de comandos
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath != "" {
		return *configPath
	}

	// Comprobar variable de entorno
	if envPath := os.Getenv("CONFIG_FILE"); envPath != "" {
		return envPath
	}

	// Usar valor por defecto
	return DefaultConfigFile
}

// createDefaultConfig crea un archivo de configuración por defecto
func createDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    DefaultPort,
			HostKey: DefaultHostKey,
		},
	}
}

// LoadConfig carga la configuración desde el archivo
func (m *Manager) LoadConfig() error {
	// Verificar si el archivo existe
	if _, err := os.Stat(m.configFile); os.IsNotExist(err) {
		// Crear configuración por defecto
		m.config = createDefaultConfig()
		// Guardar configuración por defecto
		if err := m.SaveConfig(); err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
		return nil
	}

	// Leer archivo de configuración
	data, err := os.ReadFile(m.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Decodificar JSON
	m.config = &Config{}
	if err := json.Unmarshal(data, m.config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validar configuración
	if err := m.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

// SaveConfig guarda la configuración en el archivo
func (m *Manager) SaveConfig() error {
	// Crear directorio si no existe
	dir := filepath.Dir(m.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Codificar JSON
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %v", err)
	}

	// Escribir archivo
	if err := os.WriteFile(m.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// validateConfig valida la configuración
func (m *Manager) validateConfig() error {
	// Validar puerto
	if m.config.Server.Port <= 0 || m.config.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", m.config.Server.Port)
	}

	// Validar host key
	if m.config.Server.HostKey == "" {
		return fmt.Errorf("host key file path is required")
	}

	// Validar endpoints
	if m.config.Endpoints != nil {
		// Validar default endpoint
		if m.config.Server.DefaultEndpoint != "" {
			if _, exists := m.config.Endpoints[m.config.Server.DefaultEndpoint]; !exists {
				return fmt.Errorf("default endpoint '%s' not found", m.config.Server.DefaultEndpoint)
			}
		}

		// Validar cada endpoint
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

// GetConfig devuelve la configuración actual
func (m *Manager) GetConfig() *Config {
	return m.config
}

// GetEndpoint obtiene la configuración de un endpoint específico
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
