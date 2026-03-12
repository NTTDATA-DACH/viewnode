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
