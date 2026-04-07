package ctx

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	"viewnode/cmd/config"
)

const kubeConfigFixture = `apiVersion: v1
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
- name: staging-cluster
  context:
    cluster: cluster-a
    user: user-a
current-context: dev-cluster
`

func setupConfig(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	kubeconfigPath := filepath.Join(tmpDir, "config")
	require.NoError(t, os.WriteFile(kubeconfigPath, []byte(kubeConfigFixture), 0o600))
	t.Setenv("KUBECONFIG", kubeconfigPath)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("kubeconfig", kubeconfigPath, "")
	require.NoError(t, cmd.Flags().Set("kubeconfig", kubeconfigPath))

	_, err := config.Initialize(cmd)
	require.NoError(t, err)

	return kubeconfigPath
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

func newStaticClientConfig(t *testing.T, kubeconfig string) clientcmd.ClientConfig {
	t.Helper()

	rawConfig, err := clientcmd.Load([]byte(kubeconfig))
	require.NoError(t, err)

	return clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{})
}

func TestGetCurrentRunE(t *testing.T) {
	setupConfig(t)

	output := captureStdout(t, func() {
		err := getCurrent.RunE(getCurrent, nil)
		require.NoError(t, err)
	})

	require.Equal(t, "dev-cluster\n", output)
}

func TestSetCurrentRunE(t *testing.T) {
	kubeconfigPath := setupConfig(t)

	output := captureStdout(t, func() {
		err := setCurrent.RunE(setCurrent, []string{"staging-cluster"})
		require.NoError(t, err)
	})

	require.Equal(t, "current context set to staging-cluster\n", output)

	rawConfig, err := clientcmd.LoadFromFile(kubeconfigPath)
	require.NoError(t, err)
	require.Equal(t, "staging-cluster", rawConfig.CurrentContext)
}

func TestSetCurrentRunEContextNotFound(t *testing.T) {
	setupConfig(t)

	err := setCurrent.RunE(setCurrent, []string{"missing-cluster"})
	require.EqualError(t, err, `kubernetes context "missing-cluster" not found`)
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

func TestListCmdRunE(t *testing.T) {
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
	testListCmd.Flags().StringVarP(&contextFilter, "filter", "f", "", "show only contexts according to filter")
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
