package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google_auth"
	"github.com/labstack/echo/v4"
)

func (s *server) PostSignupCallback(ctx echo.Context) error {
	res, err := google_auth.PostSignupCallback(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
