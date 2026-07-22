package cmd

import (
	"fmt"

	"github.com/karastoyanov/nowcontrol/internal/auth"
	"github.com/karastoyanov/nowcontrol/internal/client"
	"github.com/spf13/viper"
)

// newClient builds a ServiceNow client from the configured instance/username
// and the password stored in the OS credential store.
func newClient() (*client.Client, error) {
	instance := viper.GetString("instance")
	username := viper.GetString("username")

	if instance == "" {
		return nil, fmt.Errorf("missing --instance (or set it in config/env)")
	}
	if username == "" {
		return nil, fmt.Errorf("missing --username (or set it in config/env)")
	}

	password, err := auth.GetPassword(instance, username)
	if err != nil {
		return nil, fmt.Errorf("no stored credentials for %s (%s): run `nowctl auth login`", instance, username)
	}

	return client.New(instance, username, password), nil
}
