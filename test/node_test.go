package test

import (
	"testing"
	"viewnode/srv"
)

func TestNodeLoadAndFilter(t *testing.T) {
	var api MockApi
	nf := srv.NodeFilter{
		SearchText: NodeName1,
		Api:        api,
	}
	vns, err := nf.LoadAndFilter(nil)
	if err != nil {
		t.Error(err.Error())
	}
	const expectedNoNodes = 2
	if len(vns) != expectedNoNodes {
		t.Fatalf("Loading and filtering of nodes was not correct. Got: %d, expected: %d nodes.", len(vns), expectedNoNodes)
	}
	if vns[1].Name != NodeName1 {
		t.Errorf("Loading and filtering of nodes was not correct for the node name. Got: '%s', expected: '%s'", vns[0].Name, NodeName1)
	}
}
