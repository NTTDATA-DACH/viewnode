package srv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeFilter_LoadAndFilter(t *testing.T) {
	var api MockApi
	nf := NodeFilter{
		SearchText:  NodeName1,
		Api:         api,
		WithMetrics: true,
	}
	vns, err := nf.LoadAndFilter(nil)
	require.NoError(t, err)

	const expectedNoNodes = 2
	require.Equal(t, expectedNoNodes, len(vns), "loading and filtering of nodes was not correct; got: %d, expected: %d nodes.", len(vns), expectedNoNodes)
	require.Equal(t, NodeName1, vns[1].Name, "loading and filtering of nodes was not correct for the node name; got: '%s', expected: '%s'", vns[1].Name, NodeName1)
}

func TestNodeFilter_NodeListError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: NodeListError,
	}
	nf := NodeFilter{
		Api: api,
	}
	_, err := nf.LoadAndFilter(nil)
	require.EqualError(t, err, "node list error")
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
	require.ErrorIs(t, err, ErrMetricsServerNotInstalled)
}
