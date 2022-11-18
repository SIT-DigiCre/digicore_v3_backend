package tool

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
)

func GetTool(ctx echo.Context) (api.ResGetTool, *response.Error) {
	return api.ResGetTool{DiscordUrl: env.DiscordJoinUrl}, nil
}
