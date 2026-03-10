package cmd

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseNamespaces(t *testing.T) {
	namespaces := parseNamespaces(" first,second, first , ,third ")

	assert.DeepEqual(t, namespaces, []string{"first", "second", "third"})
}
