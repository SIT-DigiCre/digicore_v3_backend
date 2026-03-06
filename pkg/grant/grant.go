package grant

import (
	"sort"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
)

// claimToGrants は claim と利用可能 grant の対応を一元管理する。
var claimToGrants = map[string][]api.ResGetUserMeGrantsGrants{
	"infra": {
		api.GroupAdmin,
		api.ForceCheckout,
		api.MailBroadcast,
		api.ActivityRecordEditOther,
	},
	"account": {
		api.BudgetAdmin,
		api.PaymentAdmin,
	},
}

// ResolveFromClaims は claim 一覧から重複を除いた grant 一覧を返す。
func ResolveFromClaims(claims []string) []api.ResGetUserMeGrantsGrants {
	grantMap := map[api.ResGetUserMeGrantsGrants]struct{}{}

	for _, claim := range claims {
		grants, ok := claimToGrants[claim]
		if !ok {
			continue
		}
		for _, g := range grants {
			grantMap[g] = struct{}{}
		}
	}

	resolved := make([]api.ResGetUserMeGrantsGrants, 0, len(grantMap))
	for g := range grantMap {
		resolved = append(resolved, g)
	}

	sort.Slice(resolved, func(i, j int) bool {
		return resolved[i] < resolved[j]
	})

	return resolved
}
