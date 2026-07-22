/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var recordGetFields string

// recordGetCmd represents the record get command
var recordGetCmd = &cobra.Command{
	Use:   "get <table> <sys_id>",
	Short: "Get a single record from a ServiceNow table by sys_id",
	Long: `Fetches a single record via GET /api/now/table/{table}/{sys_id}.

Use --fields to limit the response to specific columns
(sysparm_fields), e.g. --fields "number,short_description,state".`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		table, sysID := args[0], args[1]

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

		out, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	recordCmd.AddCommand(recordGetCmd)

	recordGetCmd.Flags().StringVar(&recordGetFields, "fields", "", "comma-separated list of fields to return")
}
