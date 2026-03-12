package ns

import "testing"

import "github.com/stretchr/testify/require"

func TestNsCmdHelpText(t *testing.T) {
	require.Equal(t, "ns", NsCmd.Use)
	require.Equal(t, "Manage Kubernetes namespaces", NsCmd.Short)
	require.Equal(t, "Manage Kubernetes namespaces and namespace-related operations.", NsCmd.Long)
}
