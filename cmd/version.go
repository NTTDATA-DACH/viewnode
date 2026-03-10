package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var version string
var commit string
var buildTime = time.Now().Format(time.RFC3339)

func currentBuildYear() int {
	t, err := time.Parse(time.RFC3339, buildTime)
	if err != nil {
		return time.Now().Year()
	}
	return t.Year()
}

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
		fmt.Printf("viewnode %s (%s) © %d NTT DATA Deutschland SE, Adam Boczek | source: https://github.com/NTTDATA-DACH/viewnode\n", version, commit, currentBuildYear())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
