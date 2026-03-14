package ns

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"viewnode/cmd/config"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestNsCmdHelpText(t *testing.T) {
	require.Equal(t, "ns", NsCmd.Use)
	require.Equal(t, "Manage Kubernetes namespaces", NsCmd.Short)
	require.Equal(t, "Manage Kubernetes namespaces and namespace-related operations.", NsCmd.Long)
}

func TestNsCmdRegistersCurrentStateSubcommands(t *testing.T) {
	getCurrentCmd, _, err := NsCmd.Find([]string{"get-current"})
	require.NoError(t, err)
	require.Same(t, getCurrent, getCurrentCmd)
	require.Equal(t, []string{"gc", "getcurrent"}, getCurrentCmd.Aliases)

	setCurrentCmd, _, err := NsCmd.Find([]string{"set-current"})
	require.NoError(t, err)
	require.Same(t, setCurrent, setCurrentCmd)
	require.Equal(t, []string{"sc", "setcurrent"}, setCurrentCmd.Aliases)
}

const kubeConfigFixtureWithNamespace = `apiVersion: v1
kind: Config
preferences: {}
users:
- name: user-a
  user:
    token: abc
clusters:
- name: cluster-a
  cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:8080
contexts:
- name: dev-cluster
  context:
    cluster: cluster-a
    user: user-a
    namespace: team-a
current-context: dev-cluster
`

const kubeConfigFixtureWithoutNamespace = `apiVersion: v1
kind: Config
preferences: {}
users:
- name: user-a
  user:
    token: abc
clusters:
- name: cluster-a
  cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:8080
contexts:
- name: dev-cluster
  context:
    cluster: cluster-a
    user: user-a
current-context: dev-cluster
`

func writeKubeconfigFixture(t *testing.T, contents string) string {
	t.Helper()

	tmpDir := t.TempDir()
	kubeconfigPath := filepath.Join(tmpDir, "config")
	require.NoError(t, os.WriteFile(kubeconfigPath, []byte(contents), 0o600))

	return kubeconfigPath
}

func captureNsStdout(t *testing.T, fn func()) string {
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

func TestGetCurrentRunE(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
	})

	var capturedKubeconfig string
	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixtureWithNamespace))
	require.NoError(t, err)

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		flag := cmd.Flags().Lookup("kubeconfig")
		require.NotNil(t, flag)
		capturedKubeconfig = flag.Value.String()
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return *rawConfig, nil
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	nsCmd := &cobra.Command{Use: "ns"}
	testGetCurrentCmd := &cobra.Command{Use: "get-current", RunE: getCurrent.RunE}
	nsCmd.AddCommand(testGetCurrentCmd)
	rootCmd.AddCommand(nsCmd)

	output := captureNsStdout(t, func() {
		rootCmd.SetArgs([]string{"--kubeconfig", "/tmp/test-kubeconfig", "ns", "get-current"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "/tmp/test-kubeconfig", capturedKubeconfig)
	require.Equal(t, "team-a\n", output)
}

func TestGetCurrentRunEDefaultNamespaceFallback(t *testing.T) {
	kubeconfigPath := writeKubeconfigFixture(t, kubeConfigFixtureWithoutNamespace)
	t.Setenv("KUBECONFIG", kubeconfigPath)

	cmd := &cobra.Command{Use: "get-current"}
	cmd.Flags().String("kubeconfig", kubeconfigPath, "")
	require.NoError(t, cmd.Flags().Set("kubeconfig", kubeconfigPath))

	output := captureNsStdout(t, func() {
		err := getCurrent.RunE(cmd, nil)
		require.NoError(t, err)
	})

	require.Equal(t, "default\n", output)
}

func TestGetCurrentRunEReturnsRawConfigError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
	})

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return clientcmdapi.Config{}, errors.New("boom")
	}

	err := getCurrent.RunE(getCurrent, nil)

	require.EqualError(t, err, "getting kubernetes raw config failed (boom)")
}

func TestGetCurrentRunEReturnsInitializeError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
	})

	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return nil, errors.New("boom")
	}

	err := getCurrent.RunE(getCurrent, nil)

	require.EqualError(t, err, "boom")
}
