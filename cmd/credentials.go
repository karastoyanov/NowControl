/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// credentialKey identifies a stored password by instance/username pair,
// mirroring how aliases and other config maps are keyed in this file.
func credentialKey(instance, username string) string {
	return fmt.Sprintf("%s|%s", instance, username)
}

// storeCredential saves a password under "credentials" in the config file
// and locks the file down to owner-only permissions, since it now holds a
// plaintext secret rather than just instance/username/aliases.
func storeCredential(instance, username, password string) error {
	creds := viper.GetStringMapString("credentials")
	if creds == nil {
		creds = map[string]string{}
	}
	creds[credentialKey(instance, username)] = password
	viper.Set("credentials", creds)
	return writeConfigLocked()
}

// loadCredential retrieves a previously stored password.
func loadCredential(instance, username string) (string, error) {
	creds := viper.GetStringMapString("credentials")
	password, ok := creds[credentialKey(instance, username)]
	if !ok {
		return "", fmt.Errorf("no stored credentials for %s (%s): run `nowctl auth login`", instance, username)
	}
	return password, nil
}

// deleteCredential removes a previously stored password.
func deleteCredential(instance, username string) error {
	creds := viper.GetStringMapString("credentials")
	key := credentialKey(instance, username)
	if _, ok := creds[key]; !ok {
		return fmt.Errorf("no stored credentials for %s (%s)", instance, username)
	}
	delete(creds, key)
	viper.Set("credentials", creds)
	return writeConfigLocked()
}

// writeConfigLocked writes the config file and restricts its permissions
// to the owner only (0600), since it may contain plaintext passwords.
func writeConfigLocked() error {
	path := configFilePath()
	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	return os.Chmod(path, 0o600)
}
