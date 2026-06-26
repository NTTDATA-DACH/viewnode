package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"viewnode/cmd/ctx"
	"viewnode/cmd/ns"
	"viewnode/srv"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var namespace string
var allNamespacesFlag bool
var nodeFilter string
var podFilter string
var showContainersFlag bool
var containerViewTypeTreeFlag bool
var containerViewTypeBlockFlag bool
var showTimesFlag bool
var showRunningFlag bool
var showReqLimitsFlag bool
var showMetricsFlag bool
var verbosity string
var kubeconfig string
var watchInterval int
var watchEnabled bool

var activeRootCommand *cobra.Command

var runOnceFunc func(context.Context) error
var runWatchFunc func(context.Context, time.Duration, func(context.Context) error, func(context.Context, time.Duration) error) error
var sleepFunc func(context.Context, time.Duration) error
var withSignalNotifyContext = signal.NotifyContext
var handleRootCommandError = func(err error) error {
	log.Fatal(err)
	return nil
}

func parseNamespaces(value string) []string {
	if value == "" {
		return nil
	}

	rawNamespaces := strings.Split(value, ",")
	namespaces := make([]string, 0, len(rawNamespaces))
	seen := make(map[string]struct{}, len(rawNamespaces))
	for _, namespace := range rawNamespaces {
		namespace = strings.TrimSpace(namespace)
		if namespace == "" {
			continue
		}
		if _, ok := seen[namespace]; ok {
			continue
		}
		seen[namespace] = struct{}{}
		namespaces = append(namespaces, namespace)
	}
	return namespaces
}

func buildViewNodeDataConfig(allNamespaces bool, selectedNamespaces []string) srv.ViewNodeDataConfig {
	config := srv.ViewNodeDataConfig{
		ShowNamespaces:       allNamespaces || len(selectedNamespaces) > 0,
		GroupPodsByNamespace: allNamespaces || len(selectedNamespaces) > 1,
		SelectedNamespaces:   selectedNamespaces,
	}
	return config
}

var RootCmd = &cobra.Command{
	Use:   "viewnode",
	Short: "'viewnode' displays nodes with their pods and containers.",
	Long: `
The 'viewnode' displays nodes with their pods and containers.
Use --watch / -w to refresh every N seconds; when set without a value it defaults to 1 second, values must be >= 1, and invalid values fail before watch mode starts.
You can find the source code and usage documentation at GitHub: https://github.com/NTTDATA-DACH/viewnode.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		activeRootCommand = cmd
		if err := resolveWatchIntervalArgument(cmd.Flags().Changed("watch"), args); err != nil {
			return err
		}
		var err error
		watchEnabled, err = validateWatchInterval(cmd.Flags().Changed("watch"), watchInterval)
		if err != nil {
			return err
		}
		if !showContainersFlag && (showReqLimitsFlag || containerViewTypeTreeFlag || containerViewTypeBlockFlag) {
			log.Fatalln("you must not use -r (--show-requests-and-limits) or -b (--container-tree-view) flag without -c (--show-containers) flag")
		}

		ctx, stop := context.WithCancel(context.Background())
		ctx, signalStop := withSignalNotifyContext(ctx, watchSignals()...)
		defer signalStop()
		defer stop()

		if err := executeRootCommand(ctx); err != nil {
			return handleRootCommandError(err)
		}
		return nil
	},
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func executeRootCommand(ctx context.Context) error {
	if err := runOnceFunc(ctx); err != nil {
		return err
	}
	if !watchEnabled {
		return nil
	}
	return runWatchFunc(ctx, time.Duration(watchInterval)*time.Second, runOnceFunc, sleepFunc)
}

// validateWatchInterval relies on Cobra to reject non-integer input before RunE executes.
func validateWatchInterval(present bool, seconds int) (bool, error) {
	if !present {
		return false, nil
	}
	if seconds < 1 {
		return false, fmt.Errorf("invalid value for --watch: %d (must be >= 1 second)", seconds)
	}
	return true, nil
}

func resolveWatchIntervalArgument(present bool, args []string) error {
	if !present || len(args) == 0 {
		return nil
	}

	interval, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid argument %q for \"--watch\": %w", args[0], err)
	}

	watchInterval = interval
	return nil
}

func currentRootCommand() *cobra.Command {
	return activeRootCommand
}

func init() {
	activeRootCommand = RootCmd
	runOnceFunc = runOnce
	runWatchFunc = runWatch
	sleepFunc = productionSleep
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.SetFormatter(&nested.Formatter{
			ShowFullLevel: true,
			HideKeys:      true,
			FieldsOrder:   []string{"component", "category"},
		})
		if err := initLog(os.Stdout, verbosity); err != nil {
			return err
		}
		return nil
	}
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace to use; accepts comma-separated values")
	RootCmd.Flags().BoolVarP(&allNamespacesFlag, "all-namespaces", "A", false, "use all namespaces")
	RootCmd.Flags().StringVarP(&nodeFilter, "node-filter", "f", "", "show only nodes according to filter")
	RootCmd.Flags().StringVarP(&podFilter, "pod-filter", "p", "", "show only pods according to filter")
	RootCmd.Flags().BoolVarP(&showContainersFlag, "show-containers", "c", false, "show containers in pod")
	RootCmd.Flags().BoolVarP(&containerViewTypeTreeFlag, "container-tree-view", "b", false, "format containers in tree view, otherwise inline")
	RootCmd.Flags().BoolVar(&containerViewTypeBlockFlag, "container-block-view", false, "deprecated alias of --container-tree-view")
	_ = RootCmd.Flags().MarkDeprecated("container-block-view", "use --container-tree-view instead")
	RootCmd.Flags().BoolVarP(&showReqLimitsFlag, "show-requests-and-limits", "r", false, "show requests and limits for containers' cpu and memory (requires -c flag)")
	RootCmd.Flags().BoolVarP(&showTimesFlag, "show-pod-start-times", "t", false, "show start times of pods")
	RootCmd.Flags().BoolVar(&showRunningFlag, "show-running-only", false, "show running pods only")
	RootCmd.Flags().BoolVarP(&showMetricsFlag, "show-metrics", "m", false, "show memory footprint of nodes, pods and containers")
	RootCmd.Flags().IntVarP(&watchInterval, "watch", "w", 0, "refresh every N seconds; defaults to 1 when set without a value (must be >= 1)")
	RootCmd.Flag("watch").NoOptDefVal = "1"
	RootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "kubectl configuration file (default: ~/.kube/config or env: $KUBECONFIG)")
	RootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", log.WarnLevel.String(), "defines log level (debug, info, warn, error, fatal, panic)")

	RootCmd.AddCommand(ctx.CtxCmd, ns.NsCmd)
}

func initConfig() {
}

func initLog(out io.Writer, verbosity string) error {
	log.SetOutput(out)
	level, err := log.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	return nil
}

func getContainerViewType(flag bool) srv.ViewType {
	if flag {
		return srv.Tree
	}
	return srv.Inline
}
