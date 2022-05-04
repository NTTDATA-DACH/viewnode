package test

import (
	"testing"
	"viewnode/srv"
)

func TestPodLoadAndFilter(t *testing.T) {
	var api MockApi
	nf := srv.NodeFilter{
		Api: api,
	}
	vns, err := nf.LoadAndFilter(nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	pf := srv.PodFilter{
		Namespace:  "",
		SearchText: PodName1,
		Api:        api,
	}
	vns, err = pf.LoadAndFilter(vns)
	if err != nil {
		t.Errorf(err.Error())
	}
	const expectedNoPods = 1
	l := len(vns[1].Pods)
	if l != expectedNoPods {
		t.Errorf("Loading and filtering of pods was not correct. Got: %d, expected %d pod.", l, expectedNoPods)
	}
}
