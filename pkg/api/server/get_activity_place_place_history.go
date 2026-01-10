package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/activity"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) GetActivityPlacePlaceHistory(ctx echo.Context, place string, params api.GetActivityPlacePlaceHistoryParams) error {
	dbClient := db.Open()

	res, err := activity.GetActivityPlacePlaceHistory(ctx, &dbClient, place, params.Period, params.Date)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}

