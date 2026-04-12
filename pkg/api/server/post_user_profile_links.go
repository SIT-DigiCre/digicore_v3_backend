package server

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
)

func (s *server) PostUserProfileLinks(ctx echo.Context) error{
	dbTransactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTransactionClient.Rollback()
	
	req := api.PostUserProfileLinksJSONBody{}
	if err := ctx.Bind(&req); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: http.StatusBadRequest, Level: "Error", Message: "リクエストの形式が不正です", Log: err.Error()})
	}

	if err := validator.Validate(req); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: http.StatusBadRequest, Level: "Error", Message: "リクエストの形式が不正です", Log: err.Error()})
	}

	res, err := user.PostUserProfileLinks(ctx, &dbTransactionClient, req)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTransactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}	
