package ns

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	require.Equal(t, []string{"getcurrent"}, getCurrentCmd.Aliases)

	setCurrentCmd, _, err := NsCmd.Find([]string{"set-current"})
	require.NoError(t, err)
	require.Same(t, setCurrent, setCurrentCmd)
	require.Equal(t, []string{"setcurrent"}, setCurrentCmd.Aliases)
}
