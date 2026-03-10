package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNamespaces(t *testing.T) {
	namespaces := parseNamespaces(" first,second, first , ,third ")

	require.Equal(t, []string{"first", "second", "third"}, namespaces)
}
