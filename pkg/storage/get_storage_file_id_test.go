package storage

import (
	"testing"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/budget"
)

type fakeStorageSelectClient struct {
	selectFunc func(dest interface{}, queryPath string, params interface{}) error
}

func (f *fakeStorageSelectClient) Select(dest interface{}, queryPath string, params interface{}) error {
	if f.selectFunc != nil {
		return f.selectFunc(dest, queryPath, params)
	}
	return nil
}

func TestValidateBudgetFileAccessAllowsNonBudgetFiles(t *testing.T) {
	client := &fakeStorageSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/budget/select_budget_file_access_from_file_id.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			accesses := dest.(*[]budget.BudgetFileAccess)
			*accesses = []budget.BudgetFileAccess{}
			return nil
		},
	}

	err := validateBudgetFileAccess(client, "third-party", "file-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBudgetFileAccessDeniesThirdPartyBudgetFile(t *testing.T) {
	client := &fakeStorageSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			switch queryPath {
			case "sql/budget/select_budget_file_access_from_file_id.sql":
				accesses := dest.(*[]budget.BudgetFileAccess)
				*accesses = []budget.BudgetFileAccess{
					{BudgetId: "budget-id", ProposerUserId: "proposer-user"},
				}
			case "sql/admin/select_user_has_claim.sql":
				rows := dest.(*[]struct {
					HasClaim bool `db:"has_claim"`
				})
				*rows = []struct {
					HasClaim bool `db:"has_claim"`
				}{
					{HasClaim: false},
				}
			default:
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return nil
		},
	}

	err := validateBudgetFileAccess(client, "third-party", "file-id")
	if err == nil {
		t.Fatal("expected permission error")
	}
	if err.Code != 403 {
		t.Fatalf("expected 403, got %d", err.Code)
	}
}

func TestValidateBudgetFileAccessAllowsProposer(t *testing.T) {
	client := &fakeStorageSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/budget/select_budget_file_access_from_file_id.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			accesses := dest.(*[]budget.BudgetFileAccess)
			*accesses = []budget.BudgetFileAccess{
				{BudgetId: "budget-id", ProposerUserId: "proposer-user"},
			}
			return nil
		},
	}

	err := validateBudgetFileAccess(client, "proposer-user", "file-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBudgetFileAccessAllowsAccountClaimUser(t *testing.T) {
	client := &fakeStorageSelectClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			switch queryPath {
			case "sql/budget/select_budget_file_access_from_file_id.sql":
				accesses := dest.(*[]budget.BudgetFileAccess)
				*accesses = []budget.BudgetFileAccess{
					{BudgetId: "budget-id", ProposerUserId: "proposer-user"},
				}
			case "sql/admin/select_user_has_claim.sql":
				rows := dest.(*[]struct {
					HasClaim bool `db:"has_claim"`
				})
				*rows = []struct {
					HasClaim bool `db:"has_claim"`
				}{
					{HasClaim: true},
				}
			default:
				t.Fatalf("unexpected query path: %s", queryPath)
			}
			return nil
		},
	}

	err := validateBudgetFileAccess(client, "account-user", "file-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
