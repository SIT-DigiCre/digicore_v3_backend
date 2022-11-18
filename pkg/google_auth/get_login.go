package google_auth

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func GetLogin(ctx echo.Context) (api.ResGetLogin, *response.Error) {
	return api.ResGetLogin{Url: loginUrl}, nil
}
