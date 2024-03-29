package google_auth

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/authenticator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

func PostLoginCallback(ctx echo.Context, dbClient db.Client, requestBody api.ReqPostLoginCallback) (api.ResPostLoginCallback, *response.Error) {
	studentNumber, err := getStudentNumberfromGoogle(requestBody.Code, loginRedirectUrl)
	if err != nil {
		return api.ResPostLoginCallback{}, err
	}
	userId, err := user.IdFromStudentNumber(dbClient, studentNumber)
	if err != nil {
		return api.ResPostLoginCallback{}, err
	}
	jwt, err := authenticator.CreateToken(userId)
	if err != nil {
		return api.ResPostLoginCallback{}, err
	}
	return api.ResPostLoginCallback{Jwt: jwt}, nil
}
