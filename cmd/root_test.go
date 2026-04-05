package cmd

import (
	"bytes"
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
	watchOn = false
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

	originalRun := RootCmd.Run
	originalPersistentPreRunE := RootCmd.PersistentPreRunE
	t.Cleanup(func() {
		RootCmd.Run = originalRun
		RootCmd.PersistentPreRunE = originalPersistentPreRunE
		resetRootCommandState()
		RootCmd.SetArgs(nil)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
	})

	var (
		capturedNamespace string
		capturedAllNS     bool
	)

	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	RootCmd.Run = func(cmd *cobra.Command, args []string) {
		capturedNamespace = namespace
		capturedAllNS = allNamespacesFlag
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

	originalRun := RootCmd.Run
	originalPersistentPreRunE := RootCmd.PersistentPreRunE
	t.Cleanup(func() {
		RootCmd.Run = originalRun
		RootCmd.PersistentPreRunE = originalPersistentPreRunE
		resetRootCommandState()
		RootCmd.SetArgs(nil)
	})

	called := false
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	RootCmd.Run = func(cmd *cobra.Command, args []string) {
		called = true
	}
	RootCmd.SetArgs([]string{})

	Execute()

	require.True(t, called)
}

func TestExecutePrintOutSendsError(t *testing.T) {
	resetRootCommandState()

	errCh := make(chan error, 1)
	executePrintOut(srv.ViewNodeData{}, errCh)

	err := <-errCh
	require.EqualError(t, err, "list of view nodes must not be null")
}

func TestExecutePrintOutNoError(t *testing.T) {
	resetRootCommandState()

	errCh := make(chan error, 1)
	executePrintOut(srv.ViewNodeData{Nodes: []srv.ViewNode{{Name: ""}}}, errCh)

	require.Empty(t, errCh)
}

func TestExecutePrintOutNodeFilterNoMatchMessage(t *testing.T) {
	resetRootCommandState()

	errCh := make(chan error, 1)
	output := captureStdout(t, func() {
		executePrintOut(srv.ViewNodeData{
			NodeFilter: "worker-a",
			Nodes: []srv.ViewNode{
				{Name: ""},
			},
		}, errCh)
	})

	require.Equal(t, "no nodes matched filter \"worker-a\"\n", output)
	require.Empty(t, errCh)
}

func TestHandleErrorsIgnoresNil(t *testing.T) {
	errCh := make(chan error, 1)
	errCh <- nil
	close(errCh)

	handleErrors(errCh)
}

func TestHandleErrorsExitsOnError(t *testing.T) {
	originalExitFunc := log.StandardLogger().ExitFunc
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
	})

	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}

	errCh := make(chan error, 1)
	errCh <- errors.New("boom")
	close(errCh)

	require.PanicsWithValue(t, 1, func() {
		handleErrors(errCh)
	})
}

func TestHandleLoadAndFilterErrorScopedEOFIncludesProxyHintAndOriginalError(t *testing.T) {
	originalExitFunc := log.StandardLogger().ExitFunc
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
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
	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}

	err := srv.DecorateError(errors.New("Get \"https://cluster.example/api/v1/nodes\": EOF"))

	require.PanicsWithValue(t, 1, func() {
		handleLoadAndFilterError(err, "node")
	})
	require.Contains(t, output.String(), "loading and filtering of nodes failed; proxy configuration may be the cause")
	require.Contains(t, output.String(), "Get \"https://cluster.example/api/v1/nodes\": EOF")
}

func TestHandleLoadAndFilterErrorScopedEOFRemainsDeterministicAcrossRepeatedCalls(t *testing.T) {
	originalExitFunc := log.StandardLogger().ExitFunc
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
		log.SetOutput(originalOutput)
		log.SetFormatter(originalFormatter)
	})

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		DisableSorting:         true,
		DisableQuote:           true,
	})
	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}

	err := srv.DecorateError(errors.New("Get \"https://cluster.example/api/v1/nodes\": EOF"))
	expectedMessage := "level=fatal msg=loading and filtering of nodes failed; proxy configuration may be the cause: scoped eof: Get \"https://cluster.example/api/v1/nodes\": EOF\n"

	for range 2 {
		var output bytes.Buffer
		log.SetOutput(&output)

		require.PanicsWithValue(t, 1, func() {
			handleLoadAndFilterError(err, "node")
		})
		require.Equal(t, expectedMessage, output.String())
	}
}

func TestHandleLoadAndFilterErrorUnauthorizedKeepsExistingFatalMessage(t *testing.T) {
	originalExitFunc := log.StandardLogger().ExitFunc
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
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
	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}

	err := srv.DecorateError(errors.New("Unauthorized"))

	require.PanicsWithValue(t, 1, func() {
		handleLoadAndFilterError(err, "node")
	})
	require.Contains(t, output.String(), "you are not authorized; please login to the cloud/cluster before continuing")
	require.NotContains(t, output.String(), "proxy configuration may be the cause")
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

	require.True(t, handleLoadAndFilterError(err, "node"))
	require.Contains(t, output.String(), "access to the node API is forbidden; node names will be extracted from the pod specification if possible")
	require.NotContains(t, output.String(), "proxy configuration may be the cause")
}

func TestHandleLoadAndFilterErrorGenericFailureKeepsExistingFatalMessage(t *testing.T) {
	originalExitFunc := log.StandardLogger().ExitFunc
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
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
	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}

	err := errors.New("dial tcp 10.0.0.1:443: i/o timeout")

	require.PanicsWithValue(t, 1, func() {
		handleLoadAndFilterError(err, "pod")
	})
	require.Contains(t, output.String(), "loading and filtering of pods failed due to: dial tcp 10.0.0.1:443: i/o timeout")
	require.NotContains(t, output.String(), "proxy configuration may be the cause")
}

func TestHandleLoadAndFilterErrorOutOfScopeEOFKeepsExistingFatalMessage(t *testing.T) {
	originalExitFunc := log.StandardLogger().ExitFunc
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() {
		log.StandardLogger().ExitFunc = originalExitFunc
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
	log.StandardLogger().ExitFunc = func(code int) {
		panic(code)
	}

	err := errors.New("Post \"https://cluster.example/api/v1/nodes\": EOF")

	require.PanicsWithValue(t, 1, func() {
		handleLoadAndFilterError(err, "node")
	})
	require.Contains(t, output.String(), "loading and filtering of nodes failed due to: Post \"https://cluster.example/api/v1/nodes\": EOF")
	require.NotContains(t, output.String(), "proxy configuration may be the cause")
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
