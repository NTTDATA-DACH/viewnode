package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
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
var containerViewTypeBlockFlag bool
var showTimesFlag bool
var showRunningFlag bool
var showReqLimitsFlag bool
var showMetricsFlag bool
var verbosity string
var kubeconfig string
var watchOn bool

var rootCmd = &cobra.Command{
	Use:   "viewnode",
	Short: "'viewnode' displays nodes with their pods and containers.",
	Long: `
The 'viewnode' displays nodes with their pods and containers.
You can find the source code and usage documentation at GitHub: https://github.com/NTTDATA-DACH/viewnode.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !showContainersFlag && (showReqLimitsFlag || containerViewTypeBlockFlag) {
			log.Fatalln("you must not use -r (--show-requests-and-limits) or -b (--container-block-view) flag without -c (--show-containers) flag")
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
	cobra.CheckErr(rootCmd.Execute())
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
	setup := srv.Setup{KubeCfgPath: kubeconfig}
	err := setup.Initialize()
	if err != nil {
		errCh <- fmt.Errorf("init setup failed (%w)", err)
		return srv.ViewNodeData{}
	}
	if namespace != "" {
		setup.Namespace = namespace
	}
	if allNamespacesFlag {
		setup.Namespace = ""
	}
	api := srv.KubernetesApi{
		Setup: &setup,
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
	var vns []srv.ViewNode
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
	vnd.Config.ShowContainers = showContainersFlag
	vnd.Config.ShowTimes = showTimesFlag
	vnd.Config.ShowReqLimits = showReqLimitsFlag
	vnd.Config.ShowMetrics = showMetricsFlag
	vnd.Config.ContainerViewType = getContainerViewType(containerViewTypeBlockFlag)

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
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace to use")
	rootCmd.Flags().BoolVarP(&allNamespacesFlag, "all-namespaces", "A", false, "use all namespaces")
	rootCmd.Flags().StringVarP(&nodeFilter, "node-filter", "f", "", "show only nodes according to filter")
	rootCmd.Flags().StringVarP(&podFilter, "pod-filter", "p", "", "show only pods according to filter")
	rootCmd.Flags().BoolVarP(&showContainersFlag, "show-containers", "c", false, "show containers in pod")
	rootCmd.Flags().BoolVarP(&containerViewTypeBlockFlag, "container-block-view", "b", false, "format view of containers as a text block, otherwise inline")
	rootCmd.Flags().BoolVarP(&showReqLimitsFlag, "show-requests-and-limits", "r", false, "show requests and limits for containers' cpu and memory (requires -c flag)")
	rootCmd.Flags().BoolVarP(&showTimesFlag, "show-pod-start-times", "t", false, "show start times of pods")
	rootCmd.Flags().BoolVar(&showRunningFlag, "show-running-only", false, "show running pods only")
	rootCmd.Flags().BoolVarP(&showMetricsFlag, "show-metrics", "m", false, "show memory footprint of nodes, pods and containers")
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", log.WarnLevel.String(), "defines log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "kubectl configuration file (default: ~/.kube/config or env: $KUBECONFIG)")
	rootCmd.PersistentFlags().BoolVarP(&watchOn, "watch", "w", false, "executes the command every second so that changes can be observed")
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
		return srv.Block
	}
	return srv.Inline
}
