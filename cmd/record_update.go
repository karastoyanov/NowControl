/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	recordUpdateFields []string
	recordUpdateData   string
	recordUpdateFile   string
	recordUpdateFormat *string
	recordUpdateOutput *string
)

// recordUpdateCmd represents the record update command
var recordUpdateCmd = &cobra.Command{
	Use:   "update <table> <sys_id>",
	Short: "Update fields on an existing ServiceNow record",
	Long: `Patches a record via PATCH /api/now/table/{table}/{sys_id}, updating
only the fields supplied.

Supply fields with --field key=value (repeatable), --data '<json>', and/or
--file <path.json>. All three can be combined; later sources override
earlier ones in this order: --file, then --data, then --field.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		table, sysID := resolveTable(args[0]), args[1]
		if sysID == "" {
			return fmt.Errorf("sys_id must not be empty")
		}

		fields, err := resolveFields(recordUpdateData, recordUpdateFile, recordUpdateFields)
		if err != nil {
			return err
		}

		c, err := newClient()
		if err != nil {
			return err
		}

		record, err := c.Update(table, sysID, fields)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Updated %s record %s\n", table, sysID)

		return exportRecords(*recordUpdateFormat, *recordUpdateOutput, []map[string]any{record})
	},
}

func init() {
	recordCmd.AddCommand(recordUpdateCmd)

	recordUpdateCmd.Flags().StringArrayVar(&recordUpdateFields, "field", nil, "field=value (repeatable)")
	recordUpdateCmd.Flags().StringVar(&recordUpdateData, "data", "", "inline JSON object of fields")
	recordUpdateCmd.Flags().StringVar(&recordUpdateFile, "file", "", "path to a JSON file of fields")
	recordUpdateFormat, recordUpdateOutput = addExportFlags(recordUpdateCmd)
}
