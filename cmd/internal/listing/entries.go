package listing

import "sort"

type entryState struct {
	Name     string
	IsActive bool
}

type ContextEntry struct {
	Name     string
	IsActive bool
}

type NamespaceEntry struct {
	Name     string
	IsActive bool
}

func PrepareContextEntries(names []string, activeName string) []ContextEntry {
	prepared := prepareEntries(names, activeName)
	entries := make([]ContextEntry, 0, len(prepared))
	for _, entry := range prepared {
		entries = append(entries, ContextEntry{
			Name:     entry.Name,
			IsActive: entry.IsActive,
		})
	}
	return entries
}

func PrepareNamespaceEntries(names []string, activeName string) []NamespaceEntry {
	prepared := prepareEntries(names, activeName)
	entries := make([]NamespaceEntry, 0, len(prepared))
	for _, entry := range prepared {
		entries = append(entries, NamespaceEntry{
			Name:     entry.Name,
			IsActive: entry.IsActive,
		})
	}
	return entries
}

func prepareEntries(names []string, activeName string) []entryState {
	sortedNames := append([]string(nil), names...)
	sort.Strings(sortedNames)

	entries := make([]entryState, 0, len(sortedNames))
	for _, name := range sortedNames {
		entries = append(entries, entryState{
			Name:     name,
			IsActive: name == activeName,
		})
	}

	return entries
}
