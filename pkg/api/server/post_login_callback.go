package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google_auth"
	"github.com/labstack/echo/v4"
)

func (s *server) PostLoginCallback(ctx echo.Context) error {
	var requestBody api.ReqPostLoginCallback
	ctx.Bind(&requestBody)
	err := validator.Check(requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	res, err := google_auth.PostLoginCallback(ctx, db.DB, requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}