/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var recordDeleteYes bool

// recordDeleteCmd represents the record delete command
var recordDeleteCmd = &cobra.Command{
	Use:   "delete <table> <sys_id>",
	Short: "Delete a record from a ServiceNow table",
	Long: `Deletes a single record via DELETE /api/now/table/{table}/{sys_id}.

This is a destructive, irreversible operation. Prompts for confirmation
unless --yes is passed.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		table, sysID := resolveTable(args[0]), args[1]
		if sysID == "" {
			return fmt.Errorf("sys_id must not be empty")
		}

		if !recordDeleteYes {
			ok, err := confirm(fmt.Sprintf("Delete %s record %s?", table, sysID))
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintln(os.Stderr, "Aborted.")
				return nil
			}
		}

		c, err := newClient()
		if err != nil {
			return err
		}

		if err := c.Delete(table, sysID); err != nil {
			return err
		}

		fmt.Printf("Deleted %s record %s\n", table, sysID)
		return nil
	},
}

func init() {
	recordCmd.AddCommand(recordDeleteCmd)

	recordDeleteCmd.Flags().BoolVarP(&recordDeleteYes, "yes", "y", false, "skip confirmation prompt")
}
