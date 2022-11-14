package work

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutWorkWorkWorkId(ctx echo.Context, dbClient db.TransactionClient, workId string, requestBody api.ReqPutWorkWorkWorkId) (api.ResGetWorkWorkWorkId, *response.Error) {
	return api.ResGetWorkWorkWorkId{}, nil
}
