package discord

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func GetUserMeDiscord(ctx echo.Context) (api.ResGetUserMeDiscord, *response.Error) {
	return api.ResGetUserMeDiscord{Url: loginUrl}, nil
}
