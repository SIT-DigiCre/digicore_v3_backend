package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/work"
	"github.com/labstack/echo/v4"
)

func (s *server) GetWorkTagTagId(ctx echo.Context, tagId string) error {
	dbClient := db.Open()

	res, err := work.GetWorkTagTagId(ctx, &dbClient, tagId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
