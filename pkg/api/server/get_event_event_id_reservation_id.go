package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/labstack/echo/v4"
)

func (s *server) GetEventEventIdReservationId(ctx echo.Context, eventId string, reservationId string) error {
	dbClient := db.Open()

	res, err := event.GetEventEventIdReservationId(ctx, &dbClient, eventId, reservationId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
