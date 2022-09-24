package users

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMe(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMe) (api.ResGetUserMe, *response.Error) {

	return api.ResGetUserMe{}, nil
}
