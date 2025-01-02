package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type KeyManager struct {
	KeyPath string
}

func NewKeyManager(keyPath string) *KeyManager {
	// Create a new KeyManager with the given key path
	return &KeyManager{
		KeyPath: keyPath,
	}
}

func (km *KeyManager) GenerateHostKey() error {
	// Generate a new host key
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyPEM, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}

	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)

	return os.WriteFile(km.KeyPath, privateKeyBytes, 0600)
}

func (km *KeyManager) LoadHostKey() (ssh.Signer, error) {
	// Load the host key from file
	privateBytes, err := os.ReadFile(km.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read host key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host key: %v", err)
	}

	return private, nil
}

func (km *KeyManager) EnsureHostKey() (ssh.Signer, error) {
	// Ensure the host key exists, generate if not
	if _, err := os.Stat(km.KeyPath); os.IsNotExist(err) {
		if err := km.GenerateHostKey(); err != nil {
			return nil, fmt.Errorf("failed to generate host key: %v", err)
		}
	}

	return km.LoadHostKey()
}
