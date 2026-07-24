/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage table name aliases",
	Long: `Commands for defining short aliases for ServiceNow table names, e.g.
"computer" for "cmdb_ci_computer", stored under "aliases" in the config
file (default $HOME/.nowctl.yaml).

Aliases are accepted anywhere a <table> argument is expected -- in
"table list", "record get/create/update/delete" -- and are offered
alongside real table names in shell completion.`,
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}

// resolveTable expands table if it's a configured alias, otherwise returns
// it unchanged so real table names keep working normally alongside
// aliases.
func resolveTable(table string) string {
	aliases := viper.GetStringMapString("aliases")
	target, ok := aliases[table]
	if !ok {
		return table
	}
	fmt.Fprintf(os.Stderr, "Resolved alias %q -> %q\n", table, target)
	return target
}
