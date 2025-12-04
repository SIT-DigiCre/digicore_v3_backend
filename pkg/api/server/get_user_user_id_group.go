package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/labstack/echo/v4"
)

// GetUserUserIdGroup は、指定されたユーザーが参加しているグループ一覧を返す
func (s *server) GetUserUserIdGroup(ctx echo.Context, userId string) error {
	dbClient := db.Open()

	res, err := group.GetUserUserIdGroup(ctx, &dbClient, userId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}


