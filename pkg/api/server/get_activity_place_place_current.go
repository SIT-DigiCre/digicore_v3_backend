package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/activity"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) GetActivityPlacePlaceCurrent(ctx echo.Context, place string) error {
	dbClient := db.Open()

	res, err := activity.GetActivityPlacePlaceCurrent(ctx, &dbClient, place)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}

