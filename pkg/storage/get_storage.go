package storage

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func GetStorage(ctx echo.Context) (api.ResGetStorage, *response.Error) {
	return api.ResGetStorage{}, nil
}
