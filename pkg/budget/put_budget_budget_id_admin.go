package budget

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutBudgetBudgetIdAdmin(ctx echo.Context, dbClient db.TransactionClient, budgetId string, requestBody api.ReqPutBudgetBudgetIdAdmin) (api.ResGetBudgetBudgetId, *response.Error) {
	return GetBudgetBudgetId(ctx, dbClient, "")
}
