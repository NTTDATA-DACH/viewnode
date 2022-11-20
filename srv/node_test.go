package srv

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestNode_LoadAndFilter(t *testing.T) {
	var api MockApi
	nf := NodeFilter{
		SearchText: NodeName1,
		Api:        api,
	}
	vns, err := nf.LoadAndFilter(nil)
	assert.NilError(t, err)

	const expectedNoNodes = 2
	assert.Equal(t, expectedNoNodes, len(vns), "loading and filtering of nodes was not correct; got: %d, expected: %d nodes.", len(vns), expectedNoNodes)
	assert.Equal(t, NodeName1, vns[1].Name, "loading and filtering of nodes was not correct for the node name; got: '%s', expected: '%s'", vns[1].Name, NodeName1)
}
