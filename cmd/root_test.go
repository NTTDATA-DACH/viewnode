package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	"viewnode/srv"
)

func TestParseNamespaces(t *testing.T) {
	namespaces := parseNamespaces(" first,second, first , ,third ")

	require.Equal(t, []string{"first", "second", "third"}, namespaces)
}

func TestBuildViewNodeDataConfigScopedNamespaceSelection(t *testing.T) {
	config := buildViewNodeDataConfig(false, []string{"team-a"})

	require.True(t, config.ShowNamespaces)
	require.False(t, config.GroupPodsByNamespace)
	require.Equal(t, []string{"team-a"}, config.SelectedNamespaces)
}

func TestBuildViewNodeDataConfigSingleNamespaceRemainsFlatForScopedRendering(t *testing.T) {
	config := buildViewNodeDataConfig(false, parseNamespaces(" team-a, team-a "))

	require.True(t, config.ShowNamespaces)
	require.False(t, config.GroupPodsByNamespace)
	require.Equal(t, []string{"team-a"}, config.SelectedNamespaces)
}

func TestBuildViewNodeDataConfigMultiNamespaceSelection(t *testing.T) {
	config := buildViewNodeDataConfig(false, []string{"team-a", "team-b"})

	require.True(t, config.ShowNamespaces)
	require.True(t, config.GroupPodsByNamespace)
	require.Equal(t, []string{"team-a", "team-b"}, config.SelectedNamespaces)
}

func TestBuildViewNodeDataConfigAllNamespacesIgnoresScopedSelectionCount(t *testing.T) {
	config := buildViewNodeDataConfig(true, nil)

	require.True(t, config.ShowNamespaces)
	require.True(t, config.GroupPodsByNamespace)
	require.Nil(t, config.SelectedNamespaces)
}

func TestBuildViewNodeDataConfigUsesParsedNamespacesForGroupedScopedRendering(t *testing.T) {
	selectedNamespaces := parseNamespaces(" team-a,team-b, team-a ")
	config := buildViewNodeDataConfig(false, selectedNamespaces)

	require.True(t, config.ShowNamespaces)
	require.True(t, config.GroupPodsByNamespace)
	require.Equal(t, []string{"team-a", "team-b"}, config.SelectedNamespaces)
}

func resetRootCommandState() {
	namespace = ""
	allNamespacesFlag = false
	nodeFilter = ""
	podFilter = ""
	showContainersFlag = false
	containerViewTypeTreeFlag = false
	containerViewTypeBlockFlag = false
	showTimesFlag = false
	showRunningFlag = false
	showReqLimitsFlag = false
	showMetricsFlag = false
	verbosity = log.WarnLevel.String()
	kubeconfig = ""
	watchInterval = 0
	watchEnabled = false
	activeRootCommand = RootCmd
	runOnceFunc = runOnce
	runWatchFunc = runWatch
	sleepFunc = productionSleep
	handleRootCommandError = func(err error) error {
		return err
	}
	resetFlagSet(RootCmd.Flags())
	resetFlagSet(RootCmd.PersistentFlags())
}

func resetFlagSet(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	fn()

	require.NoError(t, w.Close())
	os.Stdout = originalStdout

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	return buf.String()
}

func TestRootCmdHelpIncludesNamespaceDescription(t *testing.T) {
	resetRootCommandState()

	var output bytes.Buffer
	RootCmd.SetOut(&output)
	RootCmd.SetErr(&output)
	RootCmd.SetArgs([]string{"--help"})

	err := RootCmd.Execute()

	require.NoError(t, err)
	require.Contains(t, output.String(), "namespace to use; accepts comma-separated values")
}

func TestRootCmdHelpIncludesNamespaceCommand(t *testing.T) {
	resetRootCommandState()

	var output bytes.Buffer
	RootCmd.SetOut(&output)
	RootCmd.SetErr(&output)
	RootCmd.SetArgs([]string{"--help"})

	err := RootCmd.Execute()

	require.NoError(t, err)
	require.Contains(t, output.String(), "  ns          Manage Kubernetes namespaces")
}

func TestRootCmdExecuteParsesNamespaceFlag(t *testing.T) {
	resetRootCommandState()

	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
	})

	var (
		capturedNamespace string
		capturedAllNS     bool
	)

	runOnceFunc = func(_ context.Context) error {
		capturedNamespace = namespace
		capturedAllNS = allNamespacesFlag
		return nil
	}

	RootCmd.SetArgs([]string{"--namespace", "first,second"})

	err := RootCmd.Execute()

	require.NoError(t, err)
	require.Equal(t, "first,second", capturedNamespace)
	require.False(t, capturedAllNS)
}

func TestRootCmdRejectsRequestsAndLimitsWithoutContainers(t *testing.T) {
	resetRootCommandState()

	originalExitFunc := log.StandardLogger().ExitFunc
	originalPersistentPreRunE := RootCmd.PersistentPreRunE
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
		RootCmd.PersistentPreRunE = originalPersistentPreRunE
		resetRootCommandState()
		RootCmd.SetArgs(nil)
	})

	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	RootCmd.SetArgs([]string{"--show-requests-and-limits"})

	require.PanicsWithValue(t, 1, func() {
		_ = RootCmd.Execute()
	})
}

func TestExecuteRunsRootCommand(t *testing.T) {
	resetRootCommandState()

	t.Cleanup(func() {
		resetRootCommandState()
		RootCmd.SetArgs(nil)
	})

	called := false
	runOnceFunc = func(_ context.Context) error {
		called = true
		return nil
	}
	RootCmd.SetArgs([]string{})

	Execute()

	require.True(t, called)
}

func TestExecutePrintOutSendsError(t *testing.T) {
	resetRootCommandState()

	err := executePrintOut(srv.ViewNodeData{})
	require.EqualError(t, err, "list of view nodes must not be null")
}

func TestExecutePrintOutNoError(t *testing.T) {
	resetRootCommandState()

	err := executePrintOut(srv.ViewNodeData{Nodes: []srv.ViewNode{{Name: ""}}})
	require.NoError(t, err)
}

func TestExecutePrintOutNodeFilterNoMatchMessage(t *testing.T) {
	resetRootCommandState()

	output := captureStdout(t, func() {
		err := executePrintOut(srv.ViewNodeData{
			NodeFilter: "worker-a",
			Nodes: []srv.ViewNode{
				{Name: ""},
			},
		})
		require.NoError(t, err)
	})

	require.Equal(t, "no nodes matched filter \"worker-a\"\n", output)
}

func TestHandleLoadAndFilterErrorScopedEOFIncludesProxyHintAndOriginalError(t *testing.T) {
	err := srv.DecorateError(errors.New("Get \"https://cluster.example/api/v1/nodes\": EOF"))

	handled, returnedErr := handleLoadAndFilterError(err, "node")
	require.False(t, handled)
	require.ErrorContains(t, returnedErr, "loading and filtering of nodes failed; proxy configuration may be the cause")
	require.ErrorContains(t, returnedErr, "Get \"https://cluster.example/api/v1/nodes\": EOF")
}

func TestHandleLoadAndFilterErrorScopedEOFRemainsDeterministicAcrossRepeatedCalls(t *testing.T) {
	err := srv.DecorateError(errors.New("Get \"https://cluster.example/api/v1/nodes\": EOF"))
	expectedMessage := "loading and filtering of nodes failed; proxy configuration may be the cause: scoped eof: Get \"https://cluster.example/api/v1/nodes\": EOF"

	for range 2 {
		handled, returnedErr := handleLoadAndFilterError(err, "node")
		require.False(t, handled)
		require.EqualError(t, returnedErr, expectedMessage)
	}
}

func TestHandleLoadAndFilterErrorUnauthorizedKeepsExistingFatalMessage(t *testing.T) {
	err := srv.DecorateError(errors.New("Unauthorized"))

	handled, returnedErr := handleLoadAndFilterError(err, "node")
	require.False(t, handled)
	require.EqualError(t, returnedErr, "you are not authorized; please login to the cloud/cluster before continuing")
}

func TestHandleLoadAndFilterErrorForbiddenKeepsWarningFallback(t *testing.T) {
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
		log.SetFormatter(originalFormatter)
	})

	var output bytes.Buffer
	log.SetOutput(&output)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		DisableSorting:         true,
		DisableQuote:           true,
	})

	err := srv.DecorateError(errors.New("nodes is forbidden: User \"alice\" cannot list resource \"nodes\""))

	handled, returnedErr := handleLoadAndFilterError(err, "node")
	require.True(t, handled)
	require.NoError(t, returnedErr)
	require.Contains(t, output.String(), "access to the node API is forbidden; node names will be extracted from the pod specification if possible")
	require.NotContains(t, output.String(), "proxy configuration may be the cause")
}

func TestHandleLoadAndFilterErrorGenericFailureKeepsExistingFatalMessage(t *testing.T) {
	err := errors.New("dial tcp 10.0.0.1:443: i/o timeout")

	handled, returnedErr := handleLoadAndFilterError(err, "pod")
	require.False(t, handled)
	require.EqualError(t, returnedErr, "loading and filtering of pods failed due to: dial tcp 10.0.0.1:443: i/o timeout")
}

func TestHandleLoadAndFilterErrorOutOfScopeEOFKeepsExistingFatalMessage(t *testing.T) {
	err := errors.New("Post \"https://cluster.example/api/v1/nodes\": EOF")

	handled, returnedErr := handleLoadAndFilterError(err, "node")
	require.False(t, handled)
	require.EqualError(t, returnedErr, "loading and filtering of nodes failed due to: Post \"https://cluster.example/api/v1/nodes\": EOF")
}

func TestInitLogSetsLoggerOutputAndLevel(t *testing.T) {
	var out bytes.Buffer

	err := initLog(&out, log.DebugLevel.String())

	require.NoError(t, err)
	require.Equal(t, log.DebugLevel, log.GetLevel())

	log.Debug("debug message")
	require.Contains(t, out.String(), "debug message")
}

func TestInitLogRejectsUnknownLevel(t *testing.T) {
	err := initLog(io.Discard, "not-a-level")

	require.Error(t, err)
}

func TestGetContainerViewType(t *testing.T) {
	require.Equal(t, srv.Tree, getContainerViewType(true))
	require.Equal(t, srv.Inline, getContainerViewType(false))
}

func TestInitConfigDoesNothing(t *testing.T) {
	initConfig()
}
