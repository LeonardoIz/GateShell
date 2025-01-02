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
    return &KeyManager{
        KeyPath: keyPath,
    }
}

func (km *KeyManager) GenerateHostKey() error {
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
    if _, err := os.Stat(km.KeyPath); os.IsNotExist(err) {
        if err := km.GenerateHostKey(); err != nil {
            return nil, fmt.Errorf("failed to generate host key: %v", err)
        }
    }

    return km.LoadHostKey()
}