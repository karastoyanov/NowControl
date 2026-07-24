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
	recordCreateFields []string
	recordCreateData   string
	recordCreateFile   string
	recordCreateFormat *string
	recordCreateOutput *string
)

// recordCreateCmd represents the record create command
var recordCreateCmd = &cobra.Command{
	Use:   "create <table>",
	Short: "Create a new record in a ServiceNow table",
	Long: `Creates a record via POST /api/now/table/{table}.

Supply fields with --field key=value (repeatable), --data '<json>', and/or
--file <path.json>. All three can be combined; later sources override
earlier ones in this order: --file, then --data, then --field.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		table := args[0]

		fields, err := resolveFields(recordCreateData, recordCreateFile, recordCreateFields)
		if err != nil {
			return err
		}

		c, err := newClient()
		if err != nil {
			return err
		}

		record, err := c.Create(table, fields)
		if err != nil {
			return err
		}

		if sysID, ok := record["sys_id"].(string); ok {
			fmt.Fprintf(os.Stderr, "Created %s record %s\n", table, sysID)
		}

		return exportRecords(*recordCreateFormat, *recordCreateOutput, []map[string]any{record})
	},
}

func init() {
	recordCmd.AddCommand(recordCreateCmd)

	recordCreateCmd.Flags().StringArrayVar(&recordCreateFields, "field", nil, "field=value (repeatable)")
	recordCreateCmd.Flags().StringVar(&recordCreateData, "data", "", "inline JSON object of fields")
	recordCreateCmd.Flags().StringVar(&recordCreateFile, "file", "", "path to a JSON file of fields")
	recordCreateFormat, recordCreateOutput = addExportFlags(recordCreateCmd)
}
