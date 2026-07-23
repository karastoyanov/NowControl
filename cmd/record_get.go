/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	recordGetFields string
	recordGetFormat *string
	recordGetOutput *string
)

// recordGetCmd represents the record get command
var recordGetCmd = &cobra.Command{
	Use:   "get <table> <sys_id>",
	Short: "Get a single record from a ServiceNow table by sys_id",
	Long: `Fetches a single record via GET /api/now/table/{table}/{sys_id}.

Use --fields to limit the response to specific columns
(sysparm_fields), e.g. --fields "number,short_description,state".

Use --format/--output to export the record as csv, xml, or xlsx instead of
printing JSON to stdout.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		table, sysID := args[0], args[1]
		if sysID == "" {
			return fmt.Errorf("sys_id must not be empty")
		}

		c, err := newClient()
		if err != nil {
			return err
		}

		params := url.Values{}
		if recordGetFields != "" {
			params.Set("sysparm_fields", recordGetFields)
		}

		record, err := c.Get(table, sysID, params)
		if err != nil {
			return err
		}

		return exportRecords(*recordGetFormat, *recordGetOutput, []map[string]any{record})
	},
}

func init() {
	recordCmd.AddCommand(recordGetCmd)

	recordGetCmd.Flags().StringVar(&recordGetFields, "fields", "", "comma-separated list of fields to return")
	recordGetFormat, recordGetOutput = addExportFlags(recordGetCmd)
}
