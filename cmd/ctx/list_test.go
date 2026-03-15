package ctx

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"viewnode/cmd/config"
)

func TestListCmdRunEUsesSharedContextEntries(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	setup := &config.Setup{
		ClientConfig: newStaticClientConfig(t, kubeConfigFixture),
	}
	var capturedKubeconfig string

	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		flag := cmd.Flags().Lookup("kubeconfig")
		require.NotNil(t, flag)
		capturedKubeconfig = flag.Value.String()
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	ctxCmd := &cobra.Command{Use: "ctx"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	ctxCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(ctxCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"--kubeconfig", "/tmp/test-kubeconfig", "ctx", "list"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "/tmp/test-kubeconfig", capturedKubeconfig)
	require.Equal(t, "[*] dev-cluster\n[ ] staging-cluster\n", output)
}

func TestListCmdRunEReturnsInitializeErrorWithoutPartialOutput(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return nil, errors.New("failed to initialize setup (config file not found at the following path: /tmp/missing-config)")
	}

	var runErr error
	output := captureStdout(t, func() {
		runErr = listCmd.RunE(listCmd, nil)
	})

	require.EqualError(t, runErr, "failed to initialize setup (config file not found at the following path: /tmp/missing-config)")
	require.Empty(t, output)
}
