package auth

import (
	"golang.org/x/crypto/ssh"
)

// PasswordAuthMethod returns an AuthMethod for password authentication
func PasswordAuthMethod(password string) ssh.AuthMethod {
	return ssh.Password(password)
}

// NoneAuthMethod returns an AuthMethod for no authentication
func NoneAuthMethod() ssh.AuthMethod {
	return ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
		return []ssh.Signer{}, nil
	})
}
