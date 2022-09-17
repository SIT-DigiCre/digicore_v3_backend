package google_auth

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/authenticator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/users"
	"github.com/labstack/echo/v4"
)

func PostSignupCallback(ctx echo.Context, db db.DBClient) (api.ResPostSignupCallback, *response.Error) {
	var req api.ReqPostSignupCallback
	ctx.Bind(&req)
	studentNumber, err := getStudentNumberfromGoogle(req.Code, signupRedirectUrl)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	err = createUser(studentNumber, db)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	userID, err := users.IDFromStudentNumber(studentNumber, db)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	jwt, err := authenticator.CreateToken(userID)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	return api.ResPostSignupCallback{Jwt: jwt}, nil
}
