package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string
var commit string

func SetVersion(v string) {
	version = v
}

func SetCommit(c string) {
	commit = c
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Plugin Version",
	Long:  "Prints out the version of the plugin and the commit hash used for the build.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("viewnode %s (%s) © 2023 NTT DATA DACH\n", version, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
