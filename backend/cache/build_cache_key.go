package cache

import (
	"sort"
	"strings"
)

func BuildCacheKey(q map[string][]string) string {
	if len(q) == 0 {
		return "all"
	}

	// Sort keys for consistent ordering
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		values := q[k]
		// Sort values for consistent ordering
		sort.Strings(values)
		parts = append(parts, k+"="+strings.Join(values, ","))
	}

	return strings.Join(parts, "&")
}
