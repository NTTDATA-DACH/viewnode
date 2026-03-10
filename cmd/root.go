package cmd

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"viewnode/cmd/config"
	"viewnode/cmd/ctx"
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
var watchOn bool

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

var RootCmd = &cobra.Command{
	Use:   "viewnode",
	Short: "'viewnode' displays nodes with their pods and containers.",
	Long: `
The 'viewnode' displays nodes with their pods and containers.
You can find the source code and usage documentation at GitHub: https://github.com/NTTDATA-DACH/viewnode.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !showContainersFlag && (showReqLimitsFlag || containerViewTypeTreeFlag || containerViewTypeBlockFlag) {
			log.Fatalln("you must not use -r (--show-requests-and-limits) or -b (--container-tree-view) flag without -c (--show-containers) flag")
		}
		stopCh := make(chan bool)
		errCh := make(chan error)
		go handleErrors(errCh)
		var wg sync.WaitGroup
		wg.Add(1)
		go schedule(&wg, stopCh, errCh)
		if !watchOn {
			close(stopCh)
		}
		wg.Wait()
	},
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func schedule(wg *sync.WaitGroup, stop <-chan bool, errCh chan<- error) {
	defer wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	vnd := executeLoadAndFilter(errCh)
	executePrintOut(vnd, errCh)
	for {
		select {
		case <-stop:
			ticker.Stop()
			return
		case <-ticker.C:
			vnd = executeLoadAndFilter(errCh)
			executePrintOut(vnd, errCh)
		}
	}
}

func executeLoadAndFilter(errCh chan<- error) srv.ViewNodeData {
	setup := config.GetConfig()
	selectedNamespaces := parseNamespaces(namespace)
	if namespace != "" {
		setup.Namespace = strings.Join(selectedNamespaces, ",")
	}
	if allNamespacesFlag {
		setup.Namespace = ""
		selectedNamespaces = nil
	}
	api := srv.KubernetesApi{
		Setup: setup,
	}
	fs := []srv.LoadAndFilter{
		srv.NodeFilter{
			SearchText:  nodeFilter,
			Api:         api,
			WithMetrics: showMetricsFlag,
		},
		srv.PodFilter{
			Namespace:   setup.Namespace,
			SearchText:  podFilter,
			Api:         api,
			RunningOnly: showRunningFlag,
			WithMetrics: showMetricsFlag,
		},
	}
	var (
		vns []srv.ViewNode
		err error
	)
	for _, f := range fs {
		log.Tracef("starting loading and filtering of %ss", f.ResourceName())
		vns, err = f.LoadAndFilter(vns)
		if err != nil {
			log.Debugf("ERROR: %s", err.Error())
			switch {
			case errors.As(err, &srv.UnauthorizedError{}):
				log.Fatalln("you are not authorized; please login to the cloud/cluster before continuing")
			case errors.As(err, &srv.NodesIsForbiddenError{}):
				log.Warnln("access to the node API is forbidden; node names will be extracted from the pod specification if possible")
				continue
			case errors.Is(err, srv.ErrMetricsServerNotInstalled):
				log.Warnf("loading of metrics for %ss failed; %s", f.ResourceName(), err.Error())
				continue
			case strings.Contains(err.Error(), "net/http: TLS handshake timeout"):
				log.Fatalf("loading and filtering of %ss failed; is the cluster up and running?", f.ResourceName())
			default:
				log.Fatalf("loading and filtering of %ss failed due to: %s", f.ResourceName(), err.Error())
			}
		}
		log.Tracef("finished loading and filtering of %ss", f.ResourceName())
	}
	vnd := srv.ViewNodeData{
		Namespace: setup.Namespace,
		Nodes:     vns,
	}
	vnd.Config.ShowNamespaces = allNamespacesFlag
	if len(selectedNamespaces) > 1 {
		vnd.Config.ShowNamespaces = true
	}
	vnd.Config.ShowContainers = showContainersFlag
	vnd.Config.ShowTimes = showTimesFlag
	vnd.Config.ShowReqLimits = showReqLimitsFlag
	vnd.Config.ShowMetrics = showMetricsFlag
	vnd.Config.ContainerViewType = getContainerViewType(containerViewTypeTreeFlag || containerViewTypeBlockFlag)

	return vnd
}

func executePrintOut(vnd srv.ViewNodeData, errCh chan<- error) {
	err := vnd.Printout(watchOn)
	if err != nil {
		errCh <- err
		return
	}
}

func handleErrors(errCh <-chan error) {
	for err := range errCh {
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func init() {
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
	RootCmd.Flags().BoolVarP(&watchOn, "watch", "w", false, "executes the command every second so that changes can be observed")
	RootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "kubectl configuration file (default: ~/.kube/config or env: $KUBECONFIG)")
	RootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", log.WarnLevel.String(), "defines log level (debug, info, warn, error, fatal, panic)")

	RootCmd.AddCommand(ctx.CtxCmd)

	_, err := config.Initialize(RootCmd)
	if err != nil {
		log.Fatal(err)
	}
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
