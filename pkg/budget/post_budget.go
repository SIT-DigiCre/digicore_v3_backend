package budget

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostBudget(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostBudget) (api.ResGetBudgetBudgetId, *response.Error) {
	return GetBudgetBudgetId(ctx, dbClient, "")
}
