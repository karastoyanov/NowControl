/*
Copyright © 2026 Aleksander Karastoyanov
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/karastoyanov/nowcontrol/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check nowctl's config, credentials, and connectivity",
	Long: `Checks the resolved config file, instance, username, and stored
credentials, then verifies connectivity and authentication against the
ServiceNow instance via GET /api/now/table/sys_user.

Exits non-zero if any check fails.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runDoctor()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

// doctorCheck is one row of `nowctl doctor` output.
type doctorCheck struct {
	name   string
	ok     bool
	detail string
}

// runDoctor runs each check in order, prints a report, and exits non-zero
// if any check failed. It exits directly rather than returning an error
// from RunE, since a failed health-check summary shouldn't also be
// followed by cobra's usage-text dump.
func runDoctor() {
	var checks []doctorCheck
	problems := 0
	fail := func(name, detail string) {
		checks = append(checks, doctorCheck{name, false, detail})
		problems++
	}
	pass := func(name, detail string) {
		checks = append(checks, doctorCheck{name, true, detail})
	}

	path := configFilePath()
	info, err := os.Stat(path)
	switch {
	case err != nil:
		pass("Config file", fmt.Sprintf("%s not found (create one with `nowctl auth login`)", path))
	case info.Mode().Perm()&0o077 != 0:
		fail("Config file", fmt.Sprintf("%s has permissions %04o, should be 0600 -- it may hold a plaintext password readable by other local users", path, info.Mode().Perm()))
	default:
		pass("Config file", fmt.Sprintf("%s (permissions %04o)", path, info.Mode().Perm()))
	}

	instance := viper.GetString("instance")
	if instance == "" {
		fail("Instance", "not set -- pass --instance, set NOWCTL_INSTANCE, or run `nowctl auth login`")
	} else {
		pass("Instance", instance)
	}

	username := viper.GetString("username")
	if username == "" {
		fail("Username", "not set -- pass --username, set NOWCTL_USERNAME, or run `nowctl auth login`")
	} else {
		pass("Username", username)
	}

	var password string
	if instance != "" && username != "" {
		password, err = loadCredential(instance, username)
		if err != nil {
			fail("Credentials", "no password stored for this instance/username -- run `nowctl auth login`")
		} else {
			pass("Credentials", "stored")
		}
	} else {
		fail("Credentials", "skipped (instance/username not resolved)")
	}

	if password != "" {
		c := client.New(instance, username, password)
		if err := c.Ping(); err != nil {
			fail("Authentication", err.Error())
		} else {
			pass("Authentication", "verified via GET sys_user")
		}
	} else {
		fail("Authentication", "skipped (missing instance/username/credentials)")
	}

	for _, c := range checks {
		status := "OK"
		if !c.ok {
			status = "FAIL"
		}
		fmt.Printf("%-16s %-4s %s\n", c.name+":", status, c.detail)
	}

	if problems > 0 {
		fmt.Printf("\n%d problem(s) found.\n", problems)
		os.Exit(1)
	}
	fmt.Println("\nAll checks passed.")
}
