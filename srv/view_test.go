package srv

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestViewNodeDataPrintoutNamespaceLabel(t *testing.T) {
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	assert.NilError(t, err)
	os.Stdout = w

	printErr := ViewNodeData{
		Namespace: "kube-system,local-path-storage",
		Nodes: []ViewNode{
			{Name: ""},
		},
	}.Printout(false)

	_ = w.Close()
	os.Stdout = originalStdout
	assert.NilError(t, printErr)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(buf.String(), "namespace(s): kube-system,local-path-storage\n"))
}
