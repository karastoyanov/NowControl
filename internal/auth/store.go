// Package auth stores and retrieves ServiceNow credentials using the OS
// credential store (Keychain on macOS, Credential Manager on Windows,
// Secret Service on Linux) via zalando/go-keyring.
package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const service = "nowctl"

func key(instance, username string) string {
	return fmt.Sprintf("%s|%s", instance, username)
}

// SetPassword stores a password for the given instance/username pair.
func SetPassword(instance, username, password string) error {
	return keyring.Set(service, key(instance, username), password)
}

// GetPassword retrieves a previously stored password.
func GetPassword(instance, username string) (string, error) {
	return keyring.Get(service, key(instance, username))
}

// DeletePassword removes a previously stored password.
func DeletePassword(instance, username string) error {
	return keyring.Delete(service, key(instance, username))
}
