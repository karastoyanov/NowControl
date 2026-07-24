/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"net/url"

	"github.com/karastoyanov/nowcontrol/internal/client"
	"github.com/spf13/cobra"
)

var (
	tableDescribeFormat *string
	tableDescribeOutput *string
)

// tableDescribeCmd represents the table describe command
var tableDescribeCmd = &cobra.Command{
	Use:   "describe <table>",
	Short: "List a table's fields, types, and reference targets",
	Long: `Describes a ServiceNow table's schema via sys_dictionary, including
fields inherited from parent tables in the class hierarchy -- e.g.
"incident" inherits number, short_description, state, etc. from "task",
and those are included alongside incident's own fields.

Field types and reference targets are shown as their display values
(e.g. "Reference" rather than a raw internal_type sys_id, "Configuration
Item" rather than the target table's sys_id).

Use --format/--output to export as csv, xml, or xlsx instead of printing
JSON to stdout, same as table list.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		table := resolveTable(args[0])

		c, err := newClient()
		if err != nil {
			return err
		}

		fields, err := describeTable(c, table)
		if err != nil {
			return err
		}

		return exportRecords(*tableDescribeFormat, *tableDescribeOutput, fields)
	},
}

func init() {
	tableCmd.AddCommand(tableDescribeCmd)
	tableDescribeFormat, tableDescribeOutput = addExportFlags(tableDescribeCmd)
}

// describeTable returns field definitions for table and all of its
// ancestor tables in the class hierarchy, deduplicated by field name so
// the most specific definition wins.
func describeTable(c *client.Client, table string) ([]map[string]any, error) {
	chain, err := tableHierarchy(c, table)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var fields []map[string]any

	for _, t := range chain {
		params := url.Values{
			"sysparm_query":         []string{fmt.Sprintf("name=%s^active=true^elementISNOTEMPTY^ORDERBYelement", t)},
			"sysparm_fields":        []string{"element,column_label,internal_type,max_length,mandatory,read_only,reference,default_value"},
			"sysparm_display_value": []string{"true"},
			"sysparm_limit":         []string{"1000"},
		}
		records, err := c.List("sys_dictionary", params)
		if err != nil {
			return nil, err
		}
		for _, r := range records {
			element, _ := r["element"].(string)
			if element == "" || seen[element] {
				continue
			}
			seen[element] = true
			fields = append(fields, r)
		}
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields found for table %q (check the table name)", table)
	}

	return fields, nil
}

// tableHierarchy resolves table and each of its ancestor tables (via
// sys_db_object.super_class), most specific first, e.g.
// ["incident", "task"] for a table extending task.
func tableHierarchy(c *client.Client, table string) ([]string, error) {
	chain := []string{table}
	current := table

	// Bounded to guard against unexpected cyclic super_class data; real
	// ServiceNow class hierarchies are never anywhere near this deep.
	for range 10 {
		params := url.Values{
			"sysparm_query":  []string{"name=" + current},
			"sysparm_fields": []string{"name,super_class"},
			"sysparm_limit":  []string{"1"},
		}
		records, err := c.List("sys_db_object", params)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			if len(chain) == 1 {
				return nil, fmt.Errorf("unknown table %q", table)
			}
			break
		}

		superClassID, ok := referenceValue(records[0]["super_class"])
		if !ok {
			break
		}

		parentParams := url.Values{
			"sysparm_query":  []string{"sys_id=" + superClassID},
			"sysparm_fields": []string{"name"},
			"sysparm_limit":  []string{"1"},
		}
		parentRecords, err := c.List("sys_db_object", parentParams)
		if err != nil {
			return nil, err
		}
		if len(parentRecords) == 0 {
			break
		}
		parentName, _ := parentRecords[0]["name"].(string)
		if parentName == "" || parentName == current {
			break
		}

		chain = append(chain, parentName)
		current = parentName
	}

	return chain, nil
}

// referenceValue extracts the sys_id from a ServiceNow reference field's
// raw (non-display-value) JSON shape -- {"link": "...", "value": "<sys_id>"}
// -- returning ("", false) if the field is empty (e.g. no parent table).
func referenceValue(v any) (string, bool) {
	m, ok := v.(map[string]any)
	if !ok {
		return "", false
	}
	value, _ := m["value"].(string)
	return value, value != ""
}
