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
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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
	originalCurrentRawConfig := currentRawConfig
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		listNamespaces = originalListNamespaces
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixtureWithNamespace))
	require.NoError(t, err)
	rawConfig.Contexts[rawConfig.CurrentContext].Namespace = "team-b"

	setup := &config.Setup{}
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
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return *rawConfig, nil
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

func TestListCmdRunEAcceptsNamespaceFilterFlags(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "long flag",
			args: []string{"--filter", "team"},
		},
		{
			name: "short flag",
			args: []string{"-f", "team"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalInitializeConfig := initializeConfig
			originalCurrentSetup := currentSetup
			originalCurrentRawConfig := currentRawConfig
			originalListNamespaces := listNamespaces
			t.Cleanup(func() {
				initializeConfig = originalInitializeConfig
				currentSetup = originalCurrentSetup
				currentRawConfig = originalCurrentRawConfig
				listNamespaces = originalListNamespaces
			})

			rawConfig, err := clientcmd.Load([]byte(kubeConfigFixtureWithNamespace))
			require.NoError(t, err)

			setup := &config.Setup{}
			initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
				return setup, nil
			}
			currentSetup = func() *config.Setup {
				return setup
			}
			currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
				require.Equal(t, setup, actualSetup)
				return *rawConfig, nil
			}
			listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
				require.Equal(t, setup, actualSetup)
				return []string{"team-a"}, nil
			}

			rootCmd := &cobra.Command{Use: "viewnode"}
			rootCmd.PersistentFlags().String("kubeconfig", "", "")
			rootCmd.Flags().StringP("node-filter", "f", "", "")
			nsCmd := &cobra.Command{Use: "ns"}
			testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
			testListCmd.Flags().StringVarP(&namespaceFilter, "filter", "f", "", "show only namespaces according to filter")
			nsCmd.AddCommand(testListCmd)
			rootCmd.AddCommand(nsCmd)

			output := captureStdout(t, func() {
				rootCmd.SetArgs(append([]string{"ns", "list"}, tc.args...))
				err := rootCmd.Execute()
				require.NoError(t, err)
			})

			require.Equal(t, "[*] team-a\n", output)
		})
	}
}

func TestListCmdRunEFiltersNamespacesCaseInsensitivelyBeforeFormatting(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		listNamespaces = originalListNamespaces
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixtureWithNamespace))
	require.NoError(t, err)
	rawConfig.Contexts[rawConfig.CurrentContext].Namespace = "kube-system"

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return *rawConfig, nil
	}
	listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
		require.Equal(t, setup, actualSetup)
		return []string{"team-a", "Kube-public", "kube-system", "team-b"}, nil
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	nsCmd := &cobra.Command{Use: "ns"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&namespaceFilter, "filter", "f", "", "show only namespaces according to filter")
	nsCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(nsCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ns", "list", "--filter", "KuBe"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "[ ] Kube-public\n[*] kube-system\n", output)
}

func TestListCmdRunETreatsEmptyFilterLikeNoFilter(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		listNamespaces = originalListNamespaces
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixtureWithNamespace))
	require.NoError(t, err)
	rawConfig.Contexts[rawConfig.CurrentContext].Namespace = "team-b"

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return *rawConfig, nil
	}
	listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
		require.Equal(t, setup, actualSetup)
		return []string{"team-c", "team-a", "team-b"}, nil
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	nsCmd := &cobra.Command{Use: "ns"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&namespaceFilter, "filter", "f", "", "show only namespaces according to filter")
	nsCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(nsCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ns", "list", "--filter="})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "[ ] team-a\n[*] team-b\n[ ] team-c\n", output)
}

func TestListCmdRunEPrintsNoMatchMessageForFilter(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		listNamespaces = originalListNamespaces
	})

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
		require.Equal(t, setup, actualSetup)
		return []string{"team-a", "team-b"}, nil
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	nsCmd := &cobra.Command{Use: "ns"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&namespaceFilter, "filter", "f", "", "show only namespaces according to filter")
	nsCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(nsCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ns", "list", "--filter", "missing"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "no namespaces matched filter \"missing\"\n", output)
}

func TestListCmdRunEDefaultNamespaceFallback(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	originalCurrentRawConfig := currentRawConfig
	originalListNamespaces := listNamespaces
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
		currentRawConfig = originalCurrentRawConfig
		listNamespaces = originalListNamespaces
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixtureWithoutNamespace))
	require.NoError(t, err)

	setup := &config.Setup{}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}
	currentRawConfig = func(actualSetup *config.Setup) (clientcmdapi.Config, error) {
		require.Equal(t, setup, actualSetup)
		return *rawConfig, nil
	}
	listNamespaces = func(ctx context.Context, actualSetup *config.Setup) ([]string, error) {
		require.Equal(t, setup, actualSetup)
		return []string{"team-b", "default", "team-a"}, nil
	}

	output := captureStdout(t, func() {
		err := listCmd.RunE(listCmd, nil)
		require.NoError(t, err)
	})

	require.Equal(t, "[*] default\n[ ] team-a\n[ ] team-b\n", output)
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

	var runErr error
	output := captureStdout(t, func() {
		runErr = listCmd.RunE(listCmd, nil)
	})

	require.EqualError(t, runErr, "getting kubernetes namespaces failed (boom)")
	require.Empty(t, output)
}

func TestListCmdRunEReturnsInitializeErrorWithoutPartialOutput(t *testing.T) {
	originalInitializeConfig := initializeConfig
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
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
