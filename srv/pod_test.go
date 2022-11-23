package srv

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestPodFilter_LoadAndFilter(t *testing.T) {
	var api MockApi
	nf := NodeFilter{
		Api:         api,
		WithMetrics: true,
	}
	vns, err := nf.LoadAndFilter(nil)
	assert.NilError(t, err)

	pf := PodFilter{
		Namespace:   "",
		SearchText:  PodName1,
		Api:         api,
		WithMetrics: false,
	}
	vns, err = pf.LoadAndFilter(vns)
	assert.NilError(t, err)

	const expectedNoPods = 1
	l := len(vns[1].Pods)
	assert.Equal(t, expectedNoPods, l, "loading and filtering of pods was not correct; got: %d, expected %d pods", l, expectedNoPods)
}

func TestPodFilter_PodListError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: PodListError,
	}
	nf := NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	assert.NilError(t, err)
	pf := PodFilter{
		Api: api,
	}
	_, err = pf.LoadAndFilter(vns)
	assert.Error(t, err, "pod list error")
}

func TestPodFilter_PodMetricsesError(t *testing.T) {
	api := MockApi{
		ApiTypeValue: PodMetricsesError,
	}
	nf := NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	assert.NilError(t, err)
	pf := PodFilter{
		WithMetrics: true,
		Api:         api,
	}
	_, err = pf.LoadAndFilter(vns)
	assert.ErrorType(t, err, ErrMetricsServerNotInstalled)
}
