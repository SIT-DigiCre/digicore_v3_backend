package grant

import (
	"testing"
)

func TestResolveFromClaims(t *testing.T) {
	tests := []struct {
		name     string
		claims   []string
		expected []string
	}{
		{
			name:   "infra のみ",
			claims: []string{"infra"},
			expected: []string{
				"activity_record_edit_other",
				"force_checkout",
				"group_admin",
				"mail_broadcast",
			},
		},
		{
			name:   "account のみ",
			claims: []string{"account"},
			expected: []string{
				"budget_admin",
				"payment_admin",
			},
		},
		{
			name:     "claim なし",
			claims:   []string{},
			expected: []string{},
		},
		{
			name:   "infra と account の両方",
			claims: []string{"infra", "account"},
			expected: []string{
				"activity_record_edit_other",
				"budget_admin",
				"force_checkout",
				"group_admin",
				"mail_broadcast",
				"payment_admin",
			},
		},
		{
			name:   "未知 claim と重複を含む",
			claims: []string{"infra", "infra", "unknown"},
			expected: []string{
				"activity_record_edit_other",
				"force_checkout",
				"group_admin",
				"mail_broadcast",
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
