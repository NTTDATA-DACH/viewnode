package ns

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"viewnode/cmd/config"
)

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

func TestListCmdRunE(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		listNamespaces = originalListNamespaces
	})

	setup := &config.Setup{Namespace: "team-b"}
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
	listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
		require.Equal(t, setup, actualSetup)
		return []string{"team-c", "team-a", "team-b"}, nil
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	nsCmd := &cobra.Command{Use: "ns"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	nsCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(nsCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"--kubeconfig", "/tmp/test-kubeconfig", "ns", "list"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "/tmp/test-kubeconfig", capturedKubeconfig)
	require.Equal(t, "[ ] team-a\n[*] team-b\n[ ] team-c\n", output)
}

func TestListCmdRunEReturnsNamespaceDiscoveryError(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		listNamespaces = originalListNamespaces
	})

	setup := &config.Setup{Namespace: "team-a"}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
		require.Equal(t, setup, actualSetup)
		return nil, errors.New("boom")
	}

	err := listCmd.RunE(listCmd, nil)

	require.EqualError(t, err, "getting kubernetes namespaces failed (boom)")
}
