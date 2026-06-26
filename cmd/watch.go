package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"viewnode/cmd/config"
	"viewnode/srv"

	log "github.com/sirupsen/logrus"
)

// runOnce refreshes command state and renders exactly one frame for one-shot and watch mode.
func runOnce(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := config.Initialize(currentRootCommand()); err != nil {
		return err
	}

	vnd, err := executeLoadAndFilter()
	if err != nil {
		return err
	}
	return executePrintOut(vnd)
}

// runWatch is entered only after the first refresh has already succeeded via runOnce in RootCmd.RunE.
func runWatch(ctx context.Context, interval time.Duration, runOnce func(context.Context) error, sleep func(context.Context, time.Duration) error) error {
	for {
		if err := sleep(ctx, interval); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		start := time.Now()
		if err := runOnce(ctx); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			renderRefreshError(start, err)
			continue
		}
		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}

func renderRefreshError(start time.Time, err error) {
	fmt.Fprint(os.Stdout, "\033[2J\033[0;0H")
	fmt.Fprintf(os.Stdout, "[%s] watch refresh failed: %s\n", start.Format(time.RFC3339), err.Error())
}

func productionSleep(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func executeLoadAndFilter() (srv.ViewNodeData, error) {
	setup := config.GetConfig()
	selectedNamespaces := parseNamespaces(namespace)
	if namespace != "" {
		setup.Namespace = strings.Join(selectedNamespaces, ",")
	}
	if allNamespacesFlag {
		setup.Namespace = ""
		selectedNamespaces = nil
	}

	api := srv.KubernetesApi{Setup: setup}
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
		filteredNodes, err := f.LoadAndFilter(vns)
		if err != nil {
			log.Debugf("ERROR: %s", err.Error())
			shouldContinue, handledErr := handleLoadAndFilterError(err, f.ResourceName())
			if handledErr != nil {
				return srv.ViewNodeData{}, handledErr
			}
			if shouldContinue {
				continue
			}
		}
		vns = filteredNodes
		log.Tracef("finished loading and filtering of %ss", f.ResourceName())
	}

	vnd := srv.ViewNodeData{
		Namespace:  setup.Namespace,
		NodeFilter: nodeFilter,
		Nodes:      vns,
	}
	vnd.Config = buildViewNodeDataConfig(allNamespacesFlag, selectedNamespaces)
	vnd.Config.ShowContainers = showContainersFlag
	vnd.Config.ShowTimes = showTimesFlag
	vnd.Config.ShowReqLimits = showReqLimitsFlag
	vnd.Config.ShowMetrics = showMetricsFlag
	vnd.Config.ContainerViewType = getContainerViewType(containerViewTypeTreeFlag || containerViewTypeBlockFlag)

	return vnd, nil
}

func executePrintOut(vnd srv.ViewNodeData) error {
	return vnd.Printout(watchEnabled)
}

func handleLoadAndFilterError(err error, resourceName string) (bool, error) {
	switch {
	case errors.As(err, &srv.UnauthorizedError{}):
		return false, errors.New("you are not authorized; please login to the cloud/cluster before continuing")
	case errors.As(err, &srv.NodesIsForbiddenError{}):
		log.Warnln("access to the node API is forbidden; node names will be extracted from the pod specification if possible")
		return true, nil
	case errors.Is(err, srv.ErrMetricsServerNotInstalled):
		log.Warnf("loading of metrics for %ss failed; %s", resourceName, err.Error())
		return true, nil
	case errors.As(err, &srv.ScopedEOFError{}):
		return false, fmt.Errorf("loading and filtering of %ss failed; proxy configuration may be the cause: %w", resourceName, err)
	case strings.Contains(err.Error(), "net/http: TLS handshake timeout"):
		return false, fmt.Errorf("loading and filtering of %ss failed; is the cluster up and running?", resourceName)
	default:
		return false, fmt.Errorf("loading and filtering of %ss failed due to: %w", resourceName, err)
	}
}
