package google_auth

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func GetSignup(ctx echo.Context) (api.ResGetSignup, *response.Error) {
	return api.ResGetSignup{Url: signupUrl}, nil
}

func PostSignupCallback(ctx echo.Context) (api.ResPostSignupCallback, *response.Error) {
	return api.ResPostSignupCallback{Jwt: "sample"}, nil
}
