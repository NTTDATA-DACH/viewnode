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
var nodeFilter string

var rootCmd = &cobra.Command{
	Use:   "kubectl-viewnode",
	Short: "kubectl-viewnode shows nodes with their pods and containers.",
	Long:  `kubectl-viewnode shows nodes with their pods and containers.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := srv.InitSetup()
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
			panic(err.Error())
		}
		pf := srv.PodFilter{
			Namespace:  setup.Namespace,
			SearchText: "",
			Api:        api,
		}
		vns, err = pf.LoadAndFilter(vns)
		vnd := srv.ViewNodeData{
			Nodes: vns,
		}
		err = vnd.Printout()
		if err != nil {
			panic(err.Error())
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
	rootCmd.Flags().StringVarP(&nodeFilter, "node-filter", "f", "", "show only nodes containing filter in node name")
}

func initConfig() {
}
