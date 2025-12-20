package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/budget"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) PutBudgetBudgetIdStatusPaid(ctx echo.Context, budgetId string) error {
	var requestBody api.ReqPutBudgetBudgetIdStatusPaid
	if err := ctx.Bind(&requestBody); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "リクエストボディの解析に失敗しました。正しい形式で送信してください", Log: err.Error()})
	}
	err := validator.Validate(requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	dbTranisactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTranisactionClient.Rollback()

	res, err := budget.PutBudgetBudgetIdStatusPaid(ctx, &dbTranisactionClient, budgetId, requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTranisactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
