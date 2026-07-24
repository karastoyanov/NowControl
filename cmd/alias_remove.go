/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasRemoveCmd represents the alias remove command
var aliasRemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove a table alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		aliases := viper.GetStringMapString("aliases")
		if _, ok := aliases[alias]; !ok {
			return fmt.Errorf("no alias named %q", alias)
		}
		delete(aliases, alias)
		viper.Set("aliases", aliases)

		if err := writeConfigLocked(); err != nil {
			return err
		}

		fmt.Printf("Removed alias %q from %s\n", alias, configFilePath())
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasRemoveCmd)
}
