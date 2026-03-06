package grant

import (
	"testing"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
)

func TestResolveFromClaims(t *testing.T) {
	tests := []struct {
		name     string
		claims   []string
		expected []api.ResGetUserMeGrantsGrants
	}{
		{
			name:   "infra のみ",
			claims: []string{"infra"},
			expected: []api.ResGetUserMeGrantsGrants{
				api.ActivityRecordEditOther,
				api.ForceCheckout,
				api.GroupAdmin,
				api.MailBroadcast,
			},
		},
		{
			name:   "account のみ",
			claims: []string{"account"},
			expected: []api.ResGetUserMeGrantsGrants{
				api.BudgetAdmin,
				api.PaymentAdmin,
			},
		},
		{
			name:     "claim なし",
			claims:   []string{},
			expected: []api.ResGetUserMeGrantsGrants{},
		},
		{
			name:   "infra と account の両方",
			claims: []string{"infra", "account"},
			expected: []api.ResGetUserMeGrantsGrants{
				api.ActivityRecordEditOther,
				api.BudgetAdmin,
				api.ForceCheckout,
				api.GroupAdmin,
				api.MailBroadcast,
				api.PaymentAdmin,
			},
		},
		{
			name:   "未知 claim と重複を含む",
			claims: []string{"infra", "infra", "unknown"},
			expected: []api.ResGetUserMeGrantsGrants{
				api.ActivityRecordEditOther,
				api.ForceCheckout,
				api.GroupAdmin,
				api.MailBroadcast,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ResolveFromClaims(tc.claims)
			if len(actual) != len(tc.expected) {
				t.Fatalf("expected %d grants, got %d", len(tc.expected), len(actual))
			}
			for i := range tc.expected {
				if actual[i] != tc.expected[i] {
					t.Fatalf("expected grants[%d]=%s, got %s", i, tc.expected[i], actual[i])
				}
			}
		})
	}
}
