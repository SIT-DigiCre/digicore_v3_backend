package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/labstack/echo/v4"
)

func (s *server) DeleteEventEventIdReservationId(ctx echo.Context, eventId string, reservationId string) error {

	dbTranisactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTranisactionClient.Rollback()

	// 管理者であるか確認
	userId := ctx.Get("user_id").(string)
	isAdmin, err := admin.CheckUserIsAdmin(&dbTranisactionClient, userId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	if !isAdmin {
		return response.ErrorResponse(ctx, &response.Error{Code: 403, Level: "Info", Message: "管理者権限が必要です", Log: "user is not admin"})
	}

	res := event.DeleteEventEventIdReservationId(ctx, &dbTranisactionClient, eventId, reservationId)
	if res != nil {
		return response.ErrorResponse(ctx, res)
	}

	err = dbTranisactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return ctx.NoContent(204)
}

func (s *server) DeleteEventEventIdReservationIdMe(ctx echo.Context, eventId string, reservationId string) error {

	dbTranisactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTranisactionClient.Rollback()

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
