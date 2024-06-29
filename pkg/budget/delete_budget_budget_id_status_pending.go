package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func DeleteBudgetBudgetIdStatusPending(ctx echo.Context, dbClient db.TransactionClient, budgetId string) (api.BlankSuccess, *response.Error) {
	userId := ctx.Get("user_id").(string)
	now_detail, err := GetBudgetBudgetId(ctx, dbClient, budgetId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	ownerUserIds := []string{now_detail.Proposer.UserId}
	if !utils.CheckUserId(ownerUserIds, userId) {
		return api.BlankSuccess{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "編集権限がありません", Log: "Permission error"}
	}
	if now_detail.Status != "pending" {
		return api.BlankSuccess{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ステータスが一致しません", Log: "Unacceptable change"}
	}
	err = deleteBudget(dbClient, budgetId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	return api.BlankSuccess{Success: true}, nil
}

func deleteBudget(dbClient db.TransactionClient, budgetId string) *response.Error {
	params := struct {
		BudgetId string `twowaysql:"budgetId"`
	}{
		BudgetId: budgetId,
	}
	_, rerr := dbClient.Exec("sql/budget/delete_budget.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
