package group

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetGroupGroupId(ctx echo.Context, dbClient db.Client, groupId string) (api.ResGetGroupGroupId, *response.Error) {
	return api.ResGetGroupGroupId{}, nil
}
