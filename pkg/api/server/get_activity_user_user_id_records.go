package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/activity"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) GetActivityUserUserIdRecords(ctx echo.Context, userId string, params api.GetActivityUserUserIdRecordsParams) error {
	dbClient := db.Open()

	res, err := activity.GetActivityUserUserIdRecords(ctx, &dbClient, userId, params.Place, params.Offset, params.Limit)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}

