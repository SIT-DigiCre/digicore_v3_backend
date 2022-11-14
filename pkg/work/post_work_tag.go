package work

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostWorkTag(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostWorkTag) (api.ResGetWorkTagTagId, *response.Error) {

	return api.ResGetWorkTagTagId{}, nil
}
