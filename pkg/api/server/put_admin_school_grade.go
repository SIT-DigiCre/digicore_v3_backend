package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) PutAdminSchoolGrade(ctx echo.Context) error {
	dbTransactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTransactionClient.Rollback()

	res, err := admin.PutAdminSchoolGrade(ctx, &dbTransactionClient)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTransactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
