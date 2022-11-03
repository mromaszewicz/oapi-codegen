package exp

import "sort"

// SortedMapKeys returns all the keys in a map of string keys in sorted
// order.
func SortedMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
