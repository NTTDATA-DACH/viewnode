package cmd

import (
	"fmt"
	"viewnode/srv"
	"viewnode/tools"

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
	Use:   "viewnode",
	Short: "'viewnode' shows nodes with their pods and containers.",
	Long: `
The 'viewnode' shows nodes with their pods and containers.
You can find the source code and usage documentation at GitHub: https://github.com/NTTDATA-EMEA/viewnode.`,
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
		fs := []srv.LoadAndFilter{
			srv.NodeFilter{
				SearchText: nodeFilter,
				Api:        api,
			},
			srv.PodFilter{
				Namespace:   setup.Namespace,
				SearchText:  podFilter,
				Api:         api,
				RunningOnly: showRunningFlag,
			},
		}
		var vns []srv.ViewNode
		for _, f := range fs {
			tools.LogDebug(fmt.Sprintf("starting loading and filtering of %ss", f.ResourceName()), debugFlag)
			vns, err = f.LoadAndFilter(vns)
			if err != nil {
				if err.Error() == "Unauthorized" {
					tools.LogErrorAndExit(errors.Wrap(err, "warning: you are NOT authorized; please login to the cloud/cluster before continuing"), debugFlag)
				}
				tools.LogErrorAndExit(errors.Wrapf(err, "error: loading and filtering of %ss failed", f.ResourceName()), debugFlag)
			}
			tools.LogDebug(fmt.Sprintf("finished loading and filtering of %ss", f.ResourceName()), debugFlag)
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
