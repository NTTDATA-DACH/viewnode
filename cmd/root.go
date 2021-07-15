package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var debugFlag bool
var namespace string
var allNamespacesFlag bool

var rootCmd = &cobra.Command{
	Use:   "kubectl-viewnode",
	Short: "kubectl-viewnode shows nodes with their pods and containers.",
	Long:  `kubectl-viewnode shows nodes with their pods and containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		lr := clientcmd.NewDefaultClientConfigLoadingRules()
		if namespace == "" && !allNamespacesFlag {
			clientCfg, err := lr.Load()
			if err != nil {
				panic(err.Error())
			}
			if debugFlag {
				spew.Dump(clientCfg)
			}
			namespace = clientCfg.Contexts[clientCfg.CurrentContext].Namespace
			if namespace == "" {
				namespace = "default"
			}
		}
		if allNamespacesFlag {
			namespace = ""
		}
		fmt.Printf("namespace: %s\n", namespace)
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
