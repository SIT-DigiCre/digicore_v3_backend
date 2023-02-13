package group

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetGroup(ctx echo.Context, dbClient db.Client, params api.GetGroupParams) (api.ResGetGroup, *response.Error) {
	return api.ResGetGroup{}, nil
}
