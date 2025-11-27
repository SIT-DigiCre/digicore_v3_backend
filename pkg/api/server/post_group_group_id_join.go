package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *server) PostGroupGroupIdJoin(ctx echo.Context, groupId string) error {
	dbTransactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer func() {
		if err := dbTransactionClient.Rollback(); err != nil {
			logrus.Errorf("トランザクションのロールバックに失敗しました: %v", err)
		}
	}()

	res, err := group.PostGroupGroupIdJoin(ctx, &dbTransactionClient, groupId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTransactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
