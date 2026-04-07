package ctx

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"viewnode/cmd/config"
)

type failingClientConfig struct {
	rawConfigErr error
}

func (f failingClientConfig) RawConfig() (clientcmdapi.Config, error) {
	return clientcmdapi.Config{}, f.rawConfigErr
}

func (f failingClientConfig) ClientConfig() (*rest.Config, error) {
	return nil, errors.New("not implemented")
}

func (f failingClientConfig) Namespace() (string, bool, error) {
	return "", false, errors.New("not implemented")
}

func (f failingClientConfig) ConfigAccess() clientcmd.ConfigAccess {
	return nil
}

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

func TestListCmdRunEReturnsRawConfigErrorWithoutPartialOutput(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	setup := &config.Setup{
		ClientConfig: failingClientConfig{rawConfigErr: errors.New("boom")},
	}

	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}

	var runErr error
	output := captureStdout(t, func() {
		runErr = listCmd.RunE(listCmd, nil)
	})

	require.EqualError(t, runErr, "getting kubernetes raw config failed (boom)")
	require.Empty(t, output)
}

func TestNormalizeContextFilter(t *testing.T) {
	require.Equal(t, "prod", normalizeContextFilter(" prod "))
	require.Equal(t, "", normalizeContextFilter("   "))
}

func TestFilterContextNames(t *testing.T) {
	contextNames := []string{"dev-cluster", "PROD-admin", "prod-cluster"}

	require.Equal(t, contextNames, filterContextNames(contextNames, ""))
	require.Equal(t, []string{"PROD-admin", "prod-cluster"}, filterContextNames(contextNames, "prod"))
}

func TestListCmdRunEAcceptsContextFilterFlags(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "long flag",
			args: []string{"--filter", "prod"},
		},
		{
			name: "short flag",
			args: []string{"-f", "prod"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalInitializeConfig := initializeConfig
			originalCurrentSetup := currentSetup
			t.Cleanup(func() {
				initializeConfig = originalInitializeConfig
				currentSetup = originalCurrentSetup
			})

			rawConfig, err := clientcmd.Load([]byte(`apiVersion: v1
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
- name: prod-admin
  context:
    cluster: cluster-a
    user: user-a
- name: dev-cluster
  context:
    cluster: cluster-a
    user: user-a
current-context: prod-admin
`))
			require.NoError(t, err)

			setup := &config.Setup{
				ClientConfig: clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{}),
			}
			initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
				return setup, nil
			}
			currentSetup = func() *config.Setup {
				return setup
			}

			rootCmd := &cobra.Command{Use: "viewnode"}
			rootCmd.PersistentFlags().String("kubeconfig", "", "")
			rootCmd.Flags().StringP("node-filter", "f", "", "")
			ctxCmd := &cobra.Command{Use: "ctx"}
			testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
			testListCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
			ctxCmd.AddCommand(testListCmd)
			rootCmd.AddCommand(ctxCmd)

			output := captureStdout(t, func() {
				rootCmd.SetArgs(append([]string{"ctx", "list"}, tc.args...))
				err := rootCmd.Execute()
				require.NoError(t, err)
			})

			require.Equal(t, "[*] prod-admin\n", output)
		})
	}
}

func TestListCmdRunEFiltersContextsCaseInsensitivelyBeforeFormatting(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	rawConfig, err := clientcmd.Load([]byte(`apiVersion: v1
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
- name: Kube-PROD
  context:
    cluster: cluster-a
    user: user-a
- name: prod-cluster
  context:
    cluster: cluster-a
    user: user-a
current-context: prod-cluster
`))
	require.NoError(t, err)

	setup := &config.Setup{
		ClientConfig: clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{}),
	}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	ctxCmd := &cobra.Command{Use: "ctx"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
	ctxCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(ctxCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ctx", "list", "--filter", "PrOd"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "[ ] Kube-PROD\n[*] prod-cluster\n", output)
}

func TestListCmdRunETreatsEmptyFilterLikeNoFilter(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixture))
	require.NoError(t, err)

	setup := &config.Setup{
		ClientConfig: clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{}),
	}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	ctxCmd := &cobra.Command{Use: "ctx"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
	ctxCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(ctxCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ctx", "list", "--filter="})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "[*] dev-cluster\n[ ] staging-cluster\n", output)
}

func TestListCmdRunEPrintsNoMatchMessageForFilteredContexts(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixture))
	require.NoError(t, err)

	setup := &config.Setup{
		ClientConfig: clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{}),
	}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	ctxCmd := &cobra.Command{Use: "ctx"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
	ctxCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(ctxCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ctx", "list", "--filter", "does-not-exist"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "no contexts matched filter \"does-not-exist\"\n", output)
}

func TestListCmdRunEOmitsActiveMarkerWhenCurrentContextIsFilteredOut(t *testing.T) {
	originalInitializeConfig := initializeConfig
	originalCurrentSetup := currentSetup
	t.Cleanup(func() {
		initializeConfig = originalInitializeConfig
		currentSetup = originalCurrentSetup
	})

	rawConfig, err := clientcmd.Load([]byte(kubeConfigFixture))
	require.NoError(t, err)

	setup := &config.Setup{
		ClientConfig: clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{}),
	}
	initializeConfig = func(cmd *cobra.Command) (*config.Setup, error) {
		return setup, nil
	}
	currentSetup = func() *config.Setup {
		return setup
	}

	rootCmd := &cobra.Command{Use: "viewnode"}
	rootCmd.PersistentFlags().String("kubeconfig", "", "")
	rootCmd.Flags().StringP("node-filter", "f", "", "")
	ctxCmd := &cobra.Command{Use: "ctx"}
	testListCmd := &cobra.Command{Use: "list", RunE: listCmd.RunE}
	testListCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
	ctxCmd.AddCommand(testListCmd)
	rootCmd.AddCommand(ctxCmd)

	output := captureStdout(t, func() {
		rootCmd.SetArgs([]string{"ctx", "list", "--filter", "staging"})
		err := rootCmd.Execute()
		require.NoError(t, err)
	})

	require.Equal(t, "[ ] staging-cluster\n", output)
}
