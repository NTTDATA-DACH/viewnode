package cmd

import (
	"fmt"
	"kubectl-viewnode/srv"
	"kubectl-viewnode/tools"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var debugFlag bool
var namespace string
var allNamespacesFlag bool
var nodeFilter string
var podFilter string
var showContainersFlag bool
var showTimesFlag bool
var showRunningFlag bool
var showReqLimitsFlag bool

var rootCmd = &cobra.Command{
	Use:   "kubectl-viewnode",
	Short: "kubectl-viewnode shows nodes with their pods and containers.",
	Long: `
The kubectl-viewnode shows nodes with their pods and containers.
You can find the source code and usage documentation at GitHub: https://github.com/NTTDATA-EMEA/kubectl-viewnode.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !showContainersFlag && showReqLimitsFlag {
			tools.LogErrorAndExit(errors.New("error -> you must not use -r flag without -c flag"), false)
		}
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
				tools.LogErrorAndExit(errors.Wrap(err, "warning -> you are NOT authorized; please login to the cloud/cluster before continuing"), debugFlag)
			}
			tools.LogErrorAndExit(errors.Wrap(err, "error -> loading of nodes failed"), debugFlag)
		}
		pf := srv.PodFilter{
			Namespace:   setup.Namespace,
			SearchText:  podFilter,
			Api:         api,
			RunningOnly: showRunningFlag,
		}
		vns, err = pf.LoadAndFilter(vns)
		if err != nil {
			tools.LogErrorAndExit(errors.Wrap(err, "error -> loading of pods failed"), debugFlag)
		}
		vnd := srv.ViewNodeData{
			Nodes: vns,
		}
		vnd.Config.ShowNamespaces = allNamespacesFlag
		vnd.Config.ShowContainers = showContainersFlag
		vnd.Config.ShowTimes = showTimesFlag
		vnd.Config.ShowReqLimits = showReqLimitsFlag
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

	rootCmd.Flags().BoolVarP(&debugFlag, "debug", "d", false, "run in debug mode (shows stack trace in case of errors)")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace to use")
	rootCmd.Flags().BoolVarP(&allNamespacesFlag, "all-namespaces", "A", false, "use all namespaces")
	rootCmd.Flags().StringVarP(&nodeFilter, "node-filter", "f", "", "show only nodes according to filter")
	rootCmd.Flags().StringVarP(&podFilter, "pod-filter", "p", "", "show only pods according to filter")
	rootCmd.Flags().BoolVarP(&showContainersFlag, "show-containers", "c", false, "show containers in pod")
	rootCmd.Flags().BoolVarP(&showReqLimitsFlag, "show-requests-and-limits", "r", false, "show requests and limits for containers' cpu and memory (requires -c flag)")
	rootCmd.Flags().BoolVarP(&showTimesFlag, "show-pod-start-times", "t", false, "show start times of pods")
	rootCmd.Flags().BoolVar(&showRunningFlag, "show-running-only", false, "show running pods only")
}

func initConfig() {
}
