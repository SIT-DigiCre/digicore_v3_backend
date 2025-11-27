package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/budget"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *server) DeleteBudgetBudgetIdStatusApprove(ctx echo.Context, budgetId string) error {
	dbTranisactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer func() {
		if err := dbTranisactionClient.Rollback(); err != nil {
			logrus.Errorf("トランザクションのロールバックに失敗しました: %v", err)
		}
	}()

	res, err := budget.DeleteBudgetBudgetIdStatusApprove(ctx, &dbTranisactionClient, budgetId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTranisactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
