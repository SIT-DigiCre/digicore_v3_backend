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
				"CLAIM_infra",
			},
		},
		{
			name:   "account のみ",
			claims: []string{"account"},
			expected: []string{
				"CLAIM_account",
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
				"CLAIM_account",
				"CLAIM_infra",
			},
		},
		{
			name:   "未知 claim と重複を含む",
			claims: []string{"infra", "infra", "unknown"},
			expected: []string{
				"CLAIM_infra",
				"CLAIM_unknown",
			},
		},
		{
			name:   "空文字 claim は無視する",
			claims: []string{"infra", ""},
			expected: []string{
				"CLAIM_infra",
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
