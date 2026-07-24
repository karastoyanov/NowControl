package cmd

import (
	"fmt"

	"github.com/karastoyanov/nowcontrol/internal/client"
	"github.com/spf13/viper"
)

// newClient builds a ServiceNow client from the configured instance/username
// and the password stored in the config file.
func newClient() (*client.Client, error) {
	instance := viper.GetString("instance")
	username := viper.GetString("username")

	if instance == "" {
		return nil, fmt.Errorf("missing --instance (or set it in config/env)")
	}
	if username == "" {
		return nil, fmt.Errorf("missing --username (or set it in config/env)")
	}

	password, err := loadCredential(instance, username)
	if err != nil {
		return nil, err
	}

	return client.New(instance, username, password), nil
}
