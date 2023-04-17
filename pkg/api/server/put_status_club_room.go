package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/status"
	"github.com/labstack/echo/v4"
)

func (s *server) PutStatusClubRoom(ctx echo.Context) error {
	var requestBody api.ReqPutStatusClubRoom
	ctx.Bind(&requestBody)
	err := validator.Validate(requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	res, err := status.PutStatusClubRoom(ctx, requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
