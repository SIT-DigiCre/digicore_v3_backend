package work

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetWorkTagTagId(ctx echo.Context, dbClient db.Client, tagId string) (api.ResGetWorkTagTagId, *response.Error) {

	return api.ResGetWorkTagTagId{}, nil
}
