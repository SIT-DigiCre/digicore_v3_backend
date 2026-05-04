package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

type budgetFilePermissionClient interface {
	Select(dest interface{}, queryPath string, params interface{}) error
}

type BudgetFileAccess struct {
	BudgetId       string `db:"budget_id"`
	ProposerUserId string `db:"proposer_user_id"`
}

func CanViewBudgetFiles(dbClient budgetFilePermissionClient, requestUserId string, proposerUserId string) (bool, *response.Error) {
	if requestUserId == proposerUserId {
		return true, nil
	}

	hasAccountClaim, err := admin.CheckUserHasClaim(dbClient, requestUserId, "account")
	if err != nil {
		return false, err
	}

	return hasAccountClaim, nil
}

func GetBudgetFileAccessFromFileId(dbClient db.Client, fileId string) (*BudgetFileAccess, *response.Error) {
	params := struct {
		FileId string `twowaysql:"fileId"`
	}{
		FileId: fileId,
	}

	budgetFileAccesses := []BudgetFileAccess{}
	err := dbClient.Select(&budgetFileAccesses, "sql/budget/select_budget_file_access_from_file_id.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ファイルの取得に失敗しました", Log: err.Error()}
	}
	if len(budgetFileAccesses) == 0 {
		return nil, nil
	}

	return &budgetFileAccesses[0], nil
}
