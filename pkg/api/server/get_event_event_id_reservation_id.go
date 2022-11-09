package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/labstack/echo/v4"
)

func (s *server) GetEventEventIDReservationID(ctx echo.Context, eventID string, reservationID string) error {
	dbClient := db.Open()

	res, err := event.GetEventEventIDReservationID(ctx, &dbClient, eventID, reservationID)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
