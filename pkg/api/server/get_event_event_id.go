package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/status"
	"github.com/labstack/echo/v4"
)

func (s *server) GetEventEventID(ctx echo.Context, eventID string) error {
	res, err := status.GetStatus(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
