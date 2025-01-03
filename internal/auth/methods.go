package auth

import (
	"golang.org/x/crypto/ssh"
)

// PasswordAuthMethod returns an AuthMethod for password authentication
func PasswordAuthMethod(password string) ssh.AuthMethod {
	return ssh.Password(password)
}
