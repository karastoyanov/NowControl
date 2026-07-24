/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nowctl",
	Short: "A command-line interface for ServiceNow",
	Long: `nowctl is a cross-platform CLI for interacting with ServiceNow: list
table records, fetch and create individual records, and manage
authentication against a ServiceNow instance.

Run "nowctl auth login" once to verify and store credentials; instance,
username, and password are then remembered in the config file, so other
commands need no flags.

Instance and username are resolved in order: --instance/--username
flags, NOWCTL_INSTANCE/NOWCTL_USERNAME env vars, the config file
(default $HOME/.nowctl.yaml), then unset.`,
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen
// once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nowctl.yaml)")
	rootCmd.PersistentFlags().String("instance", "", "ServiceNow instance URL")
	rootCmd.PersistentFlags().String("username", "", "ServiceNow username")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "print extra diagnostic info (e.g. which config file is used)")

	viper.BindPFlag("instance", rootCmd.PersistentFlags().Lookup("instance"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".nowctl")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("NOWCTL")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// SetVersionInfo sets the version string reported by `nowctl --version`,
// populated by goreleaser's ldflags at build time (defaults to "dev" for
// local `go build`).
func SetVersionInfo(version, date string) {
	rootCmd.Version = fmt.Sprintf("%s (built %s)", version, formatBuildDate(date))
}

// formatBuildDate reformats goreleaser's RFC3339 build date into dd-MM-yyyy,
// falling back to the raw value (e.g. "unknown" for local go build) if it
// doesn't parse.
func formatBuildDate(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Format("02-01-2006")
}

// configFilePath returns the config file location that a command should
// write to: the explicit --config path if given, otherwise $HOME/.nowctl.yaml.
func configFilePath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return filepath.Join(home, ".nowctl.yaml")
}
