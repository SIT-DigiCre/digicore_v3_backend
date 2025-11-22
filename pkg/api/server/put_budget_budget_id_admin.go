package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/budget"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *server) PutBudgetBudgetIdAdmin(ctx echo.Context, budgetId string) error {
	var requestBody api.ReqPutBudgetBudgetIdAdmin
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
	defer func() {
		if err := dbTranisactionClient.Rollback(); err != nil {
			logrus.Errorf("トランザクションのロールバックに失敗しました: %v", err)
		}
	}()

	res, err := budget.PutBudgetBudgetIdAdmin(ctx, &dbTranisactionClient, budgetId, requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTranisactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
