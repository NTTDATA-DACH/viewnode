package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var namespace string
var allNamespaces bool

var rootCmd = &cobra.Command{
	Use:   "kubectl-view-node",
	Short: "kubectl-view-node shows nodes with their pods and containers.",
	Long: `kubectl-view-node shows nodes with their pods and containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		lr := clientcmd.NewDefaultClientConfigLoadingRules()
		if namespace == "" && !allNamespaces {
			clientCfg, err := lr.Load()
			if err != nil {
				panic(err.Error())
			}
			namespace = clientCfg.Contexts[clientCfg.CurrentContext].Namespace
			if namespace == "" {
				namespace = "default"
			}
		}
		if allNamespaces {
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

	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace to use")
	rootCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "all namespaces to use")
}

func initConfig() {
}
