package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/work"
	"github.com/labstack/echo/v4"
)

// GetUserUserIdWork は、指定されたユーザーが作者として含まれる作品一覧を返す
func (s *server) GetUserUserIdWork(ctx echo.Context, userId string) error {
	dbClient := db.Open()

	res, err := work.GetUserUserIdWork(ctx, &dbClient, userId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}


