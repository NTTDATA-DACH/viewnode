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

func TestListCmdRunE(t *testing.T) {
	setupConfig(t)

	output := captureStdout(t, func() {
		err := listCmd.RunE(listCmd, nil)
		require.NoError(t, err)
	})

	require.Equal(t, "[*] dev-cluster\n[ ] staging-cluster\n", output)
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
