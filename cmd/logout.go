/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored ServiceNow credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		instance := viper.GetString("instance")
		username := viper.GetString("username")

		if instance == "" {
			return fmt.Errorf("missing --instance (or set it in config/env)")
		}
		if username == "" {
			return fmt.Errorf("missing --username (or set it in config/env)")
		}

		if err := deleteCredential(instance, username); err != nil {
			return fmt.Errorf("removing credentials: %w", err)
		}

		fmt.Printf("Logged out of %s (%s)\n", instance, username)
		return nil
	},
}

func init() {
	authCmd.AddCommand(logoutCmd)
}
