/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasListCmd represents the alias list command
var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured table aliases",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases := viper.GetStringMapString("aliases")
		if len(aliases) == 0 {
			fmt.Println("No aliases configured. Add one with `nowctl alias set <alias> <table>`.")
			return nil
		}

		names := make([]string, 0, len(aliases))
		for alias := range aliases {
			names = append(names, alias)
		}
		sort.Strings(names)

		for _, alias := range names {
			fmt.Printf("%s -> %s\n", alias, aliases[alias])
		}
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)
}
