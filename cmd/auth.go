/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage ServiceNow authentication",
	Long: `Commands for authenticating against a ServiceNow instance and managing
credentials stored in the OS credential store (Keychain on macOS,
Credential Manager on Windows, Secret Service on Linux).`,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
