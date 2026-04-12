package server

import (	
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

func (s *server) DeleteUserProfileLinks(ctx echo.Context, param *api.DeleteUserProfileLinksParams) error {
	id := param.ID.String()
	dbTransactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTransactionClient.Rollback()
	
	res, err := user.DeleteUserProfileLinks(ctx, &dbTransactionClient, id)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTransactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
