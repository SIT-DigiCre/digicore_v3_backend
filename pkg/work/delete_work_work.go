package work

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteWorkWorkWorkId(ctx echo.Context, dbClient db.TransactionClient, workId string) (api.BlankSuccess, *response.Error) {

	return api.BlankSuccess{}, nil
}
