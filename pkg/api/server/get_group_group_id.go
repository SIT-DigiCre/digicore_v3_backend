package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/labstack/echo/v4"
)

func (s *server) GetGroupGroupId(ctx echo.Context, groupId string) error {
	dbClient := db.Open()

	res, err := group.GetGroupGroupId(ctx, &dbClient, groupId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
