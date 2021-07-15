package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"kubectl-viewnode/srv"
)

var debugFlag bool
var namespace string
var allNamespacesFlag bool

var rootCmd = &cobra.Command{
	Use:   "kubectl-viewnode",
	Short: "kubectl-viewnode shows nodes with their pods and containers.",
	Long:  `kubectl-viewnode shows nodes with their pods and containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := srv.GetCurrentNamespaceAndClientset()
		if err != nil {
			panic(err.Error())
		}
		if debugFlag {
			spew.Dump(setup)
		}
		if namespace != "" {
			setup.Namespace = namespace
		}
		if allNamespacesFlag {
			setup.Namespace = ""
		}
		if !allNamespacesFlag {
			fmt.Printf("namespace: %s\n", setup.Namespace)
		}
		nodes, err := srv.GetNodes(setup.Clientset)
		if err != nil {
			panic(err.Error())
		}
		vnd := srv.GetViewNodeData(nodes, "")
		srv.PrintNodes(vnd)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().BoolVarP(&debugFlag, "debug", "d", false, "run in debug mode (shows more output)")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace to use")
	rootCmd.Flags().BoolVarP(&allNamespacesFlag, "all-namespaces", "A", false, "all namespaces to use")
}

func initConfig() {
}
