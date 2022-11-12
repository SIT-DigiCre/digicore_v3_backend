package storage

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostStorage(ctx echo.Context, dbTransactionClient db.TransactionClient, requestBody api.ReqPostLoginCallback) (api.ResGetStorageFileId, *response.Error) {
	return api.ResGetStorageFileId{}, nil
}
