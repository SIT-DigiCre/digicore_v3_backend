package google_auth

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/users"
	"github.com/labstack/echo/v4"
)

func GetLogin(ctx echo.Context) (api.ResGetLogin, *response.Error) {
	return api.ResGetLogin{Url: loginUrl}, nil
}

func PostLoginCallback(ctx echo.Context, db db.DBClient) (api.ResPostLoginCallback, *response.Error) {
	var req api.ReqPostSignupCallback
	ctx.Bind(&req)
	studentNumber, err := getStudentNumberfromGoogle(req.Code, loginRedirectUrl)
	if err != nil {
		return api.ResPostLoginCallback{}, err
	}
	userID, err := users.IDFromStudentNumber(studentNumber, db)
	if err != nil {
		return api.ResPostLoginCallback{}, err
	}
	jwt, err := validator.CreateToken(userID)
	if err != nil {
		return api.ResPostLoginCallback{}, err
	}
	return api.ResPostLoginCallback{Jwt: jwt}, nil
}
