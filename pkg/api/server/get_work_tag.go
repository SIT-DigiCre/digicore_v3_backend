package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/work"
	"github.com/labstack/echo/v4"
)

func (s *server) GetWorkTag(ctx echo.Context, params api.GetWorkTagParams) error {
	dbClient := db.Open()

	res, err := work.GetWorkTag(ctx, &dbClient, params)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
