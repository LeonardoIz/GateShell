package auth

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type AuthConfig struct {
	ServerVersion string
	HostKeyFile   string
	AuthMethod    string
}

type Authenticator struct {
	config     *AuthConfig
	keyManager *KeyManager
}

func NewAuthenticator(config *AuthConfig) *Authenticator {
	return &Authenticator{
		config:     config,
		keyManager: NewKeyManager(config.HostKeyFile),
	}
}

func (a *Authenticator) ConfigureServer() (*ssh.ServerConfig, error) {
	// Configure SSH server with custom settings
	config := &ssh.ServerConfig{
		ServerVersion: a.config.ServerVersion,
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// Return permissions based on provided password
			return &ssh.Permissions{
				Extensions: map[string]string{
					"password": string(password),
				},
			}, nil
		},
	}

	signer, err := a.keyManager.EnsureHostKey()
	if err != nil {
		return nil, fmt.Errorf("failed to setup host key: %v", err)
	}

	config.AddHostKey(signer)
	return config, nil
}

func (a *Authenticator) GetUpstreamConfig(username, password string) *ssh.ClientConfig {
	var authMethod ssh.AuthMethod

	switch a.config.AuthMethod {
	case "none":
		authMethod = NoneAuthMethod()
	case "password":
		authMethod = PasswordAuthMethod(password)
	default:
		authMethod = PasswordAuthMethod(password)
	}

	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}
