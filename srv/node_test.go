package srv

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNodeFilter_LoadAndFilter(t *testing.T) {
	var api MockApi
	nf := NodeFilter{
		SearchText:  NodeName1,
		Api:         api,
		WithMetrics: true,
	}
	vns, err := nf.LoadAndFilter(nil)
	assert.NilError(t, err)

	const expectedNoNodes = 2
	assert.Equal(t, expectedNoNodes, len(vns), "loading and filtering of nodes was not correct; got: %d, expected: %d nodes.", len(vns), expectedNoNodes)
	assert.Equal(t, NodeName1, vns[1].Name, "loading and filtering of nodes was not correct for the node name; got: '%s', expected: '%s'", vns[1].Name, NodeName1)
}

func TestNodeFilter_NodeListError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: NodeListError,
	}
	nf := NodeFilter{
		Api: api,
	}
	_, err := nf.LoadAndFilter(nil)
	assert.Error(t, err, "node list error")
}

func TestNodeFilter_NodeMetricsesError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: NodeMetricsesError,
	}
	nf := NodeFilter{
		WithMetrics: true,
		Api:         api,
	}
	_, err := nf.LoadAndFilter(nil)
	assert.ErrorType(t, err, ErrMetricsServerNotInstalled)
}
