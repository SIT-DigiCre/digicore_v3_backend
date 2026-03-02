package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	grade_update "github.com/SIT-DigiCre/digicore_v3_backend/pkg/grade_update"
	"github.com/labstack/echo/v4"
)

func (s *server) GetAdminGradeUpdate(ctx echo.Context) error {
	dbClient := db.Open()

	res, err := grade_update.GetAdminGradeUpdate(ctx, &dbClient)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
