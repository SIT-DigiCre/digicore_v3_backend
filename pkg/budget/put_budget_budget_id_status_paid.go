package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutBudgetBudgetIdStatusPaid(ctx echo.Context, dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdStatusPaid) (api.ResGetBudgetBudgetId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	now_detail, err := GetBudgetBudgetId(ctx, dbClient, budgetId)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	if now_detail.Proposer.UserId != userId {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "編集権限があリません", Log: "Permission error"}
	}
	if now_detail.Status != "paid" {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ステータスが一致しません", Log: "Unacceptable change"}
	}
	err = updatePaidBudget(dbClient, budgetId, requestBody)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	return GetBudgetBudgetId(ctx, dbClient, budgetId)
}

func updatePaidBudget(dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdStatusPaid) *response.Error {
	params := struct {
		BudgetId string `twowaysql:"budgetId"`
		Remark   string `twowaysql:"remark"`
	}{
		BudgetId: budgetId,
		Remark:   requestBody.Remark,
	}
	_, rerr := dbClient.Exec("sql/budget/update_paid_budget.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
