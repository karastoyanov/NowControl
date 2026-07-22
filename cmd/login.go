/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/karastoyanov/nowcontrol/internal/auth"
	"github.com/karastoyanov/nowcontrol/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against a ServiceNow instance and store credentials",
	Long: `Prompts for a password, verifies it against the ServiceNow instance,
and stores it securely in the OS credential store (Keychain on macOS,
Credential Manager on Windows, Secret Service on Linux).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instance := viper.GetString("instance")
		username := viper.GetString("username")

		if instance == "" {
			return fmt.Errorf("missing --instance (or set it in config/env)")
		}
		if username == "" {
			return fmt.Errorf("missing --username (or set it in config/env)")
		}

		password, err := readPassword()
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}

		c := client.New(instance, username, password)
		if err := c.Ping(); err != nil {
			return fmt.Errorf("could not authenticate against %s: %w", instance, err)
		}

		if err := auth.SetPassword(instance, username, password); err != nil {
			return fmt.Errorf("storing credentials: %w", err)
		}

		viper.Set("instance", instance)
		viper.Set("username", username)
		configPath := configFilePath()
		if err := viper.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("saving %s: %w", configPath, err)
		}

		fmt.Printf("Logged in to %s as %s (saved to %s)\n", instance, username, configPath)
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}

func readPassword() (string, error) {
	fmt.Print("Password: ")

	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}
