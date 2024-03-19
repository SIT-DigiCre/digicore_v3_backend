package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/budget"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) GetBudgetBudgetId(ctx echo.Context, budgetId string) error {
	dbClient := db.Open()

	res, err := budget.GetBudgetBudgetId(ctx, &dbClient, budgetId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
