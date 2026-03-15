package listing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareContextEntriesSortsNamesAndMarksActiveEntry(t *testing.T) {
	entries := PrepareContextEntries([]string{"staging-cluster", "dev-cluster", "prod-cluster"}, "prod-cluster")

	require.Equal(t, []ContextEntry{
		{Name: "dev-cluster", IsActive: false},
		{Name: "prod-cluster", IsActive: true},
		{Name: "staging-cluster", IsActive: false},
	}, entries)
}

func TestPrepareNamespaceEntriesLeavesAllEntriesInactiveWhenActiveNameMissing(t *testing.T) {
	entries := PrepareNamespaceEntries([]string{"team-c", "team-a", "team-b"}, "team-d")

	require.Equal(t, []NamespaceEntry{
		{Name: "team-a", IsActive: false},
		{Name: "team-b", IsActive: false},
		{Name: "team-c", IsActive: false},
	}, entries)
}
