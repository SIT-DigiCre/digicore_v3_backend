package grant

import (
	"sort"
)

// claimToGrants は claim と利用可能 grant の対応を一元管理する。
var claimToGrants = map[string][]string{
	"infra": {
		"group_admin",
		"force_checkout",
		"mail_broadcast",
		"activity_record_edit_other",
	},
	"account": {
		"budget_admin",
		"payment_admin",
	},
}

// ResolveFromClaims は claim 一覧から重複を除いた grant 一覧を返す。
func ResolveFromClaims(claims []string) []string {
	grantMap := map[string]struct{}{}

	for _, claim := range claims {
		grants, ok := claimToGrants[claim]
		if !ok {
			continue
		}
		for _, g := range grants {
			grantMap[g] = struct{}{}
		}
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
