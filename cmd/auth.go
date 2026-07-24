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
credentials stored in the config file (default $HOME/.nowctl.yaml),
which is locked down to owner-only read/write permissions (0600).`,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
