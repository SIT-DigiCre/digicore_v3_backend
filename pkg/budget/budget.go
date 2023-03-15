package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func addFile(dbClient db.TransactionClient, budgetId string, fileIds []string) *response.Error {
	for _, fileId := range fileIds {
		params := struct {
			BudgetId string `twowaysql:"budgetId"`
			FileId   string `twowaysql:"fileId"`
		}{
			BudgetId: budgetId,
			FileId:   fileId,
		}
		_, rerr := dbClient.Exec("sql/budget/insert_budget_file.sql", &params, false)
		if rerr != nil {
			return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
		}
	}
	return nil
}
