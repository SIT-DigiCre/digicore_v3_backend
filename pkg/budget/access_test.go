package budget

import (
	"testing"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
)

type fakeSelectClient struct {
	selectFunc func(dest interface{}, queryPath string, params interface{}) error
}

func (f *fakeSelectClient) Select(dest interface{}, queryPath string, params interface{}) error {
	if f.selectFunc != nil {
		return f.selectFunc(dest, queryPath, params)
	}
	return nil
}

func TestCanViewBudgetFilesAllowsProposer(t *testing.T) {
	allowed, err := CanViewBudgetFiles(&fakeSelectClient{}, "request-user", "request-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected proposer to be allowed")
	}
}

func TestCanViewBudgetFilesAllowsAccountClaimUser(t *testing.T) {
	client := &fakeSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/admin/select_user_has_claim.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}

			rows := dest.(*[]struct {
				HasClaim bool `db:"has_claim"`
			})
			*rows = []struct {
				HasClaim bool `db:"has_claim"`
			}{
				{HasClaim: true},
			}
			return nil
		},
	}

	allowed, err := CanViewBudgetFiles(client, "account-user", "proposer-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected account claim user to be allowed")
	}
}

func TestCanViewBudgetFilesDeniesThirdParty(t *testing.T) {
	client := &fakeSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			rows := dest.(*[]struct {
				HasClaim bool `db:"has_claim"`
			})
			*rows = []struct {
				HasClaim bool `db:"has_claim"`
			}{
				{HasClaim: false},
			}
			return nil
		},
	}

	allowed, err := CanViewBudgetFiles(client, "third-party", "proposer-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected third party to be denied")
	}
}

func TestFilterBudgetFilesForResponseHidesFilesFromThirdParty(t *testing.T) {
	client := &fakeSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			rows := dest.(*[]struct {
				HasClaim bool `db:"has_claim"`
			})
			*rows = []struct {
				HasClaim bool `db:"has_claim"`
			}{
				{HasClaim: false},
			}
			return nil
		},
	}

	res := api.ResGetBudgetBudgetId{
		Proposer: api.ResGetBudgetBudgetIdObjectProposer{UserId: "proposer-user"},
		Files: []api.ResGetBudgetBudgetIdObjectFile{
			{FileId: "file-id", Name: "receipt.pdf"},
		},
	}

	err := filterBudgetFilesForResponse(client, "third-party", &res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 0 {
		t.Fatalf("expected files to be hidden, got %d", len(res.Files))
	}
}

func TestFilterBudgetFilesForResponseKeepsFilesForProposer(t *testing.T) {
	res := api.ResGetBudgetBudgetId{
		Proposer: api.ResGetBudgetBudgetIdObjectProposer{UserId: "proposer-user"},
		Files: []api.ResGetBudgetBudgetIdObjectFile{
			{FileId: "file-id", Name: "receipt.pdf"},
		},
	}

	err := filterBudgetFilesForResponse(&fakeSelectClient{}, "proposer-user", &res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 1 {
		t.Fatalf("expected files to remain visible, got %d", len(res.Files))
	}
}
