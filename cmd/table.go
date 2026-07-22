/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// tableCmd represents the table command
var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Query ServiceNow tables",
	Long: `Commands for reading data from ServiceNow tables via the Table API
(/api/now/table/{table}), such as listing records with optional filters.`,
}

func init() {
	rootCmd.AddCommand(tableCmd)
}
