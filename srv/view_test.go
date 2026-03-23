package srv

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViewNodeDataPrintoutNamespaceLabel(t *testing.T) {
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	printErr := ViewNodeData{
		Namespace: "kube-system,local-path-storage",
		Nodes: []ViewNode{
			{Name: ""},
		},
	}.Printout(false)

	_ = w.Close()
	os.Stdout = originalStdout
	require.NoError(t, printErr)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	require.True(t, strings.Contains(buf.String(), "namespace(s): kube-system,local-path-storage\n"))
}

func TestViewNodeDataPrintoutNodeFilterNoMatchMessage(t *testing.T) {
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	printErr := ViewNodeData{
		NodeFilter: "worker-a",
		Nodes: []ViewNode{
			{Name: ""},
		},
	}.Printout(false)

	_ = w.Close()
	os.Stdout = originalStdout
	require.NoError(t, printErr)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	require.Equal(t, "no nodes matched filter \"worker-a\"\n", buf.String())
}

func TestViewNodeDataPrintoutNoNodesWithoutFilterUsesGenericMessage(t *testing.T) {
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	printErr := ViewNodeData{
		Nodes: []ViewNode{
			{Name: ""},
		},
	}.Printout(false)

	_ = w.Close()
	os.Stdout = originalStdout
	require.NoError(t, printErr)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	require.Equal(t, "no nodes to display...\n", buf.String())
}
