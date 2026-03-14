package ns

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"viewnode/cmd/config"
)

func TestSetCurrentRunE(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	originalNamespaceExists := namespaceExists
	originalPersistRawConfig := persistRawConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		namespaceExists = originalNamespaceExists
		persistRawConfig = originalPersistRawConfig
	})

	var capturedKubeconfig string
	setup := &config.Setup{KubeCfgPath: "/tmp/actual-kubeconfig"}
	rawConfig := clientcmdapi.Config{
		CurrentContext: "dev-cluster",
		Contexts: map[string]*clientcmdapi.Context{
			"dev-cluster": {Namespace: "team-a"},
			"other":       {Namespace: "team-b"},
		},
	}

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
		return rawConfig, nil
	}
	namespaceExists = func(_ context.Context, actualSetup *config.Setup, namespace string) (bool, error) {
		require.Equal(t, setup, actualSetup)
		require.Equal(t, "team-c", namespace)
		return true, nil
	}
	persistedConfig := clientcmdapi.Config{}
	persistRawConfig = func(actualSetup *config.Setup, updatedConfig clientcmdapi.Config) error {
		require.Equal(t, setup, actualSetup)
		persistedConfig = updatedConfig
		return nil
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	nsCmd := &cobra.Command{Use: "ns"}
	testSetCurrentCmd := &cobra.Command{
		Use:  "set-current",
		Args: setCurrent.Args,
		RunE: setCurrent.RunE,
	}
	nsCmd.AddCommand(testSetCurrentCmd)
	rootCmd.AddCommand(nsCmd)

	output := captureNsStdout(t, func() {
		rootCmd.SetArgs([]string{"--kubeconfig", "/tmp/test-kubeconfig", "ns", "set-current", "team-c"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "/tmp/test-kubeconfig", capturedKubeconfig)
	require.Equal(t, "team-c", persistedConfig.Contexts["dev-cluster"].Namespace)
	require.Equal(t, "team-b", persistedConfig.Contexts["other"].Namespace)
	require.Equal(t, "current namespace set to team-c\n", output)
}

func TestSetCurrentRunEReturnsMissingNamespaceError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalNamespaceExists := namespaceExists
	originalPersistRawConfig := persistRawConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		namespaceExists = originalNamespaceExists
		persistRawConfig = originalPersistRawConfig
	})

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	namespaceExists = func(_ context.Context, actualSetup *config.Setup, namespace string) (bool, error) {
		require.Equal(t, setup, actualSetup)
		require.Equal(t, "missing", namespace)
		return false, nil
	}
	persistCalled := false
	persistRawConfig = func(_ *config.Setup, _ clientcmdapi.Config) error {
		persistCalled = true
		return nil
	}

	err := setCurrent.RunE(setCurrent, []string{"missing"})

	require.EqualError(t, err, `kubernetes namespace "missing" not found`)
	require.False(t, persistCalled)
}

func TestSetCurrentRunEReturnsNamespaceValidationError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalNamespaceExists := namespaceExists
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		namespaceExists = originalNamespaceExists
	})

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	namespaceExists = func(_ context.Context, _ *config.Setup, _ string) (bool, error) {
		return false, errors.New("boom")
	}

	err := setCurrent.RunE(setCurrent, []string{"team-c"})

	require.EqualError(t, err, `validating kubernetes namespace "team-c" failed (boom)`)
}

func TestSetCurrentRunEReturnsRawConfigError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	originalNamespaceExists := namespaceExists
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		namespaceExists = originalNamespaceExists
	})

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	namespaceExists = func(_ context.Context, _ *config.Setup, _ string) (bool, error) {
		return true, nil
	}
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return clientcmdapi.Config{}, errors.New("boom")
	}

	err := setCurrent.RunE(setCurrent, []string{"team-c"})

	require.EqualError(t, err, "getting kubernetes raw config failed (boom)")
}

func TestSetCurrentRunEReturnsPersistError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	originalNamespaceExists := namespaceExists
	originalPersistRawConfig := persistRawConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		namespaceExists = originalNamespaceExists
		persistRawConfig = originalPersistRawConfig
	})

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	namespaceExists = func(_ context.Context, _ *config.Setup, _ string) (bool, error) {
		return true, nil
	}
	currentRawConfig = func(_ *config.Setup) (clientcmdapi.Config, error) {
		return clientcmdapi.Config{
			CurrentContext: "dev-cluster",
			Contexts: map[string]*clientcmdapi.Context{
				"dev-cluster": {Namespace: "team-a"},
			},
		}, nil
	}
	persistRawConfig = func(_ *config.Setup, _ clientcmdapi.Config) error {
		return errors.New("boom")
	}

	err := setCurrent.RunE(setCurrent, []string{"team-c"})

	require.EqualError(t, err, "setting kubernetes current namespace failed (boom)")
}

func TestPersistRawConfigUsesExplicitKubeconfigPath(t *testing.T) {
	kubeconfigPath := filepath.Join(t.TempDir(), "config")
	require.NoError(t, os.WriteFile(kubeconfigPath, []byte(kubeConfigFixtureWithNamespace), 0o600))

	setup := &config.Setup{KubeCfgPath: kubeconfigPath}
	require.NoError(t, setup.DetermineClientConfig())
	rawConfig, err := currentRawConfig(setup)
	require.NoError(t, err)

	rawConfig.Contexts[rawConfig.CurrentContext].Namespace = "team-c"
	require.NoError(t, persistRawConfig(setup, rawConfig))

	updatedContents, err := os.ReadFile(kubeconfigPath)
	require.NoError(t, err)
	updatedConfig, err := clientcmd.Load(updatedContents)
	require.NoError(t, err)

	require.Equal(t, "team-c", updatedConfig.Contexts[updatedConfig.CurrentContext].Namespace)
}
