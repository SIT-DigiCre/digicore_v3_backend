package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/storage"
	"github.com/labstack/echo/v4"
)

func (s *server) GetStorageMyfile(ctx echo.Context) error {
	dbClient := db.Open()

	res, err := storage.GetStorageMyfile(ctx, &dbClient)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
