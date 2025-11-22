package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *server) DeleteEventEventIdReservationIdMe(ctx echo.Context, eventId string, reservationId string) error {

	dbTranisactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer func() {
		if err := dbTranisactionClient.Rollback(); err != nil {
			logrus.Errorf("トランザクションのロールバックに失敗しました: %v", err)
		}
	}()

	res, err := event.DeleteEventEventIdReservationIdMe(ctx, &dbTranisactionClient, eventId, reservationId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTranisactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
