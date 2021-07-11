package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-view-node",
	Short: "kubectl-view-node shows nodes with their pods and containers.",
	Long: `kubectl-view-node shows nodes with their pods and containers.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
}
