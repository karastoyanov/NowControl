/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Get, create, update, and delete individual ServiceNow records",
	Long: `Commands for working with a single record on a ServiceNow table via
the Table API (/api/now/table/{table}/{sys_id}), identified by its sys_id.`,
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
