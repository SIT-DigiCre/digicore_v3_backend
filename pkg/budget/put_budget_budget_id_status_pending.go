package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func ReqPutBudgetBudgetIdStatusPending(ctx echo.Context, dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdStatusPending) (api.ResGetBudgetBudgetId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	now_detail, err := GetBudgetBudgetId(ctx, dbClient, budgetId)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	ownerUserIds := []string{now_detail.Proposer.UserId}
	if !utils.CheckUserId(ownerUserIds, userId) {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "編集権限があリません", Log: "Permission error"}
	}
	if now_detail.Status != "pending" {
		return api.ResGetBudgetBudgetId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ステータスが一致しません", Log: "Unacceptable change"}
	}
	err = updatePendingBudget(dbClient, budgetId, requestBody)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	return GetBudgetBudgetId(ctx, dbClient, budgetId)
}

func updatePendingBudget(dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdStatusPending) *response.Error {
	params := struct {
		BudgetId      string `twowaysql:"budgetId"`
		Name          string `twowaysql:"name"`
		Budget        int    `twowaysql:"budget"`
		MattermostUrl string `twowaysql:"mattermostUrl"`
		Purpose       string `twowaysql:"purpose"`
		Remark        string `twowaysql:"remark"`
	}{
		BudgetId:      budgetId,
		Name:          requestBody.Name,
		Budget:        requestBody.Budget,
		MattermostUrl: requestBody.MattermostUrl,
		Purpose:       requestBody.Purpose,
		Remark:        requestBody.Remark,
	}
	_, rerr := dbClient.Exec("sql/budget/update_pending_budget.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
