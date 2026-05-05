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

func (s *server) PostUserProfileLinks(ctx echo.Context) error {
    dbTransactionClient, opErr := db.OpenTransaction() 
   if opErr != nil {
		return response.ErrorResponse(ctx, opErr) 
     
    }
    defer dbTransactionClient.Rollback()
    
    req := api.PostUserProfileLinksJSONBody{}
    if bindErr := ctx.Bind(&req); bindErr != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: http.StatusBadRequest, Level: "Error", Message: "...", Log: bindErr.Error()})
    }

	if err := validator.Validate(req); err != nil {
		return response.ErrorResponse(ctx, err)
	}

	res, err := user.PostUserProfileLinks(ctx, &dbTransactionClient, req)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

    if comErr := dbTransactionClient.Commit(); comErr != nil { 
        return response.ErrorResponse(ctx, comErr)
    }

    return response.SuccessResponse(ctx, res)
}
