package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/storage"
	"github.com/labstack/echo/v4"
)

func (s *server) GetStorage(ctx echo.Context) error {
	res, err := storage.GetStorage(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
