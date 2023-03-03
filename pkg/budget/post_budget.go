package budget

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostBudget(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostBudget) (api.ResGetBudgetBudgetId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	budgetID, err := createBudget(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetBudgetBudgetId{}, err
	}
	return GetBudgetBudgetId(ctx, dbClient, budgetID)
}

func createBudget(dbClient db.TransactionClient, userId string, requestBody api.ReqPostBudget) (string, *response.Error) {
	status := "pending"
	if requestBody.Class == "festival" || requestBody.Class == "fixed" {
		status = "approve"
	}
	params := struct {
		Name           string `twowaysql:"name"`
		Class          string `twowaysql:"class"`
		ProposerUserId string `twowaysql:"proposerUserId"`
		Status         string `twowaysql:"status"`
	}{
		Name:           requestBody.Name,
		Class:          requestBody.Class,
		ProposerUserId: userId,
		Status:         status,
	}
	_, rerr := dbClient.Exec("sql/budget/insert_budget.sql", &params, true)
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	budgetId, rerr := dbClient.GetId()
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return budgetId, nil
}
