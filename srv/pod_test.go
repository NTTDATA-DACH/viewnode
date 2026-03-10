package srv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPodFilter_LoadAndFilter(t *testing.T) {
	var api MockApi
	nf := NodeFilter{
		Api:         api,
		WithMetrics: true,
	}
	vns, err := nf.LoadAndFilter(nil)
	require.NoError(t, err)

	pf := PodFilter{
		Namespace:   "",
		SearchText:  PodName1,
		Api:         api,
		WithMetrics: false,
	}
	vns, err = pf.LoadAndFilter(vns)
	require.NoError(t, err)

	const expectedNoPods = 1
	l := len(vns[1].Pods)
	require.Equal(t, expectedNoPods, l, "loading and filtering of pods was not correct; got: %d, expected %d pods", l, expectedNoPods)
}

func TestPodFilter_PodListError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: PodListError,
	}
	nf := NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	require.NoError(t, err)
	pf := PodFilter{
		Api: api,
	}
	_, err = pf.LoadAndFilter(vns)
	require.EqualError(t, err, "pod list error")
}

func TestPodFilter_PodMetricsesError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: PodMetricsesError,
	}
	nf := NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	require.NoError(t, err)
	pf := PodFilter{
		WithMetrics: true,
		Api:         api,
	}
	_, err = pf.LoadAndFilter(vns)
	require.ErrorIs(t, err, ErrMetricsServerNotInstalled)
}

func TestPodFilter_LoadAndFilterMultipleNamespaces(t *testing.T) {
	var api MockApi
	nf := NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	require.NoError(t, err)

	pf := PodFilter{
		Namespace: "first, second",
		Api:       api,
	}
	vns, err = pf.LoadAndFilter(vns)
	require.NoError(t, err)

	const expectedNoPods = 6
	require.Equal(t, expectedNoPods, len(vns[1].Pods), "loading and filtering of pods across namespaces was not correct; got: %d, expected %d pods", len(vns[1].Pods), expectedNoPods)
	require.Equal(t, "first", vns[1].Pods[0].Namespace)
	require.Equal(t, "second", vns[1].Pods[3].Namespace)
}
