package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutBudgetBudgetIdStatusApprove(ctx echo.Context, dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdStatusApprove) (api.ResGetBudgetBudgetId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	now_detail, err := GetBudgetBudgetId(ctx, dbClient, budgetId)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	ownerUserIds := []string{now_detail.Proposer.UserId}
	if !utils.CheckUserId(ownerUserIds, userId) {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "編集権限があリません", Log: "Permission error"}
	}
	if now_detail.Status != "approve" {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ステータスが一致しません", Log: "Unacceptable change"}
	}
	if requestBody.Bought && len(now_detail.Files) == 0 {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusBadRequest, Level: "Error", Message: "購入済みにするには領収書を添付する必要があります", Log: "Receipt not attached"}
	}
	err = updateApproveBudget(dbClient, budgetId, requestBody)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	err = addFile(dbClient, budgetId, requestBody.Files)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	return GetBudgetBudgetId(ctx, dbClient, budgetId)
}

func updateApproveBudget(dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdStatusApprove) *response.Error {
	status := "approve"
	if requestBody.Bought {
		status = "bought"
	}
	params := struct {
		BudgetId   string `twowaysql:"budgetId"`
		Status     string `twowaysql:"status"`
		Remark     string `twowaysql:"remark"`
		Settlement int    `twowaysql:"settlement"`
	}{
		BudgetId:   budgetId,
		Status:     status,
		Remark:     requestBody.Remark,
		Settlement: requestBody.Settlement,
	}
	_, rerr := dbClient.Exec("sql/budget/update_approve_budget.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
