package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutBudgetBudgetIdAdmin(ctx echo.Context, dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdAdmin) (api.ResGetBudgetBudgetId, *response.Error) {
	now_detail, err := GetBudgetBudgetId(ctx, dbClient, budgetId)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	if now_detail.Status == "pending" {
		if !(requestBody.Status == "approve" || requestBody.Status == "reject") {
			return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "更新不可能なステータスです(acceptかrejectのみ可能)", Log: "Unacceptable change"}
		}
	}
	if now_detail.Status == "bought" {
		if !(requestBody.Status == "paid") {
			return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "更新不可能なステータスです(paidのみ可能)", Log: "unacceptable change"}
		}
	}
	if requestBody.Status == "approve" {
		userId := ctx.Get("user_id").(string)
		err = setApprove(dbClient, budgetId, userId)
		if err != nil {
			return api.ResGetBudgetBudgetId{}, err
		}
	}
	err = updateStatus(dbClient, budgetId, requestBody.Status)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	return GetBudgetBudgetId(ctx, dbClient, budgetId)
}

func updateStatus(dbClient db.TransactionClient, budgetId string, status string) *response.Error {
	params := struct {
		BudgetId string `twowaysql:"budgetId"`
		Status   string `twowaysql:"status"`
	}{
		BudgetId: budgetId,
		Status:   status,
	}
	_, err := dbClient.Exec("sql/budget/update_budget_status.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}

func setApprove(dbClient db.TransactionClient, budgetId string, approverUserId string) *response.Error {
	params := struct {
		BudgetId       string `twowaysql:"budgetId"`
		ApproverUserId string `twowaysql:"approverUserId"`
	}{
		BudgetId:       budgetId,
		ApproverUserId: approverUserId,
	}
	_, err := dbClient.Exec("sql/budget/update_budget_status.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
