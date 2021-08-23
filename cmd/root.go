package cmd

import (
	"fmt"
	"kubectl-viewnode/srv"
	"kubectl-viewnode/tools"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var debugFlag bool
var namespace string
var allNamespacesFlag bool
var nodeFilter string
var podFilter string
var showContainersFlag bool

var rootCmd = &cobra.Command{
	Use:   "kubectl-viewnode",
	Short: "kubectl-viewnode shows nodes with their pods and containers.",
	Long:  `kubectl-viewnode shows nodes with their pods and containers.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := srv.InitSetup()
		if err != nil {
			tools.LogErrorAndExit(errors.Wrap(err, "error -> init setup failed"), debugFlag)
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
		api := srv.KubernetesApi{
			Setup: setup,
		}
		vf := srv.NodeFilter{
			SearchText: nodeFilter,
			Api:        api,
		}
		var vns []srv.ViewNode
		vns, err = vf.LoadAndFilter(vns)
		if err != nil {
			if err.Error() == "Unauthorized" {
				fmt.Println("warning: you are NOT authorized; please login to the cloud/cluster before continuing...")
				os.Exit(1)
			}
			tools.LogErrorAndExit(errors.Wrap(err, "error -> loading of nodes failed"), debugFlag)
		}
		pf := srv.PodFilter{
			Namespace:  setup.Namespace,
			SearchText: podFilter,
			Api:        api,
		}
		vns, err = pf.LoadAndFilter(vns)
		if err != nil {
			tools.LogErrorAndExit(errors.Wrap(err, "error -> loading of pods failed"), debugFlag)
		}
		vnd := srv.ViewNodeData{
			Nodes: vns,
		}
		vnd.Config.CanShowNamespaces = allNamespacesFlag
		vnd.Config.CanShowContainers = showContainersFlag
		err = vnd.Printout()
		if err != nil {
			tools.LogErrorAndExit(errors.Wrap(err, "error -> printing failed"), debugFlag)
		}
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
	rootCmd.Flags().BoolVarP(&allNamespacesFlag, "all-namespaces", "A", false, "use all namespaces")
	rootCmd.Flags().StringVarP(&nodeFilter, "node-filter", "f", "", "show only nodes according to filter")
	rootCmd.Flags().StringVarP(&podFilter, "pod-filter", "p", "", "show only pods according to filter")
	rootCmd.Flags().BoolVarP(&showContainersFlag, "show-containers", "c", false, "show containers in pod")
}

func initConfig() {
}
