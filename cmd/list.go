/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	listQuery  string
	listFields string
	listLimit  int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <table>",
	Short: "List records from a ServiceNow table",
	Long: `Lists records from a ServiceNow table via GET /api/now/table/{table}.

Use --query with encoded ServiceNow query syntax (sysparm_query) to filter
results, e.g. --query "active=true^priority=1".`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		table := args[0]

		c, err := newClient()
		if err != nil {
			return err
		}

		params := url.Values{}
		if listQuery != "" {
			params.Set("sysparm_query", listQuery)
		}
		if listFields != "" {
			params.Set("sysparm_fields", listFields)
		}
		if listLimit > 0 {
			params.Set("sysparm_limit", strconv.Itoa(listLimit))
		}

		records, err := c.List(table, params)
		if err != nil {
			return err
		}

		out, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	tableCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listQuery, "query", "", "encoded sysparm_query filter")
	listCmd.Flags().StringVar(&listFields, "fields", "", "comma-separated list of fields to return")
	listCmd.Flags().IntVar(&listLimit, "limit", 10, "maximum number of records to return")
}
