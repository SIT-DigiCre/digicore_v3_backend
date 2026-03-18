package grant

import (
	"fmt"
	"sort"
)

// ResolveFromClaims は claim 一覧から CLAIM_ プレフィックス付きの grant 一覧を返す。
func ResolveFromClaims(claims []string) []string {
	grantMap := map[string]struct{}{}

	for _, claim := range claims {
		if claim == "" {
			continue
		}
		grantMap[fmt.Sprintf("CLAIM_%s", claim)] = struct{}{}
	}

	resolved := make([]string, 0, len(grantMap))
	for g := range grantMap {
		resolved = append(resolved, g)
	}

	sort.Slice(resolved, func(i, j int) bool {
		return resolved[i] < resolved[j]
	})

	return resolved
}
