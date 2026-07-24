/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasSetCmd represents the alias set command
var aliasSetCmd = &cobra.Command{
	Use:   "set <alias> <table>",
	Short: "Define (or update) a table alias",
	Long: `Saves an alias in the config file so it can be used in place of a
table name in any command:

  nowctl alias set computer cmdb_ci_computer
  nowctl table list computer --limit 5

Overwrites the alias if it already exists.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias, table := args[0], args[1]

		aliases := viper.GetStringMapString("aliases")
		if aliases == nil {
			aliases = map[string]string{}
		}
		aliases[alias] = table
		viper.Set("aliases", aliases)

		if err := writeConfigLocked(); err != nil {
			return err
		}

		fmt.Printf("Alias %q -> %q saved to %s\n", alias, table, configFilePath())
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd)
}
