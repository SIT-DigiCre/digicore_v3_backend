package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

// GetAdminUser は /admin/user のハンドラ。
// 管理者向けのユーザー一覧・検索機能を提供する。
func (s *server) GetAdminUser(ctx echo.Context, params api.GetAdminUserParams) error {
	dbClient := db.Open()

	res, err := user.GetAdminUser(ctx, &dbClient, params)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
