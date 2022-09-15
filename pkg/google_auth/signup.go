package google_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/users"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func GetSignup(ctx echo.Context) (api.ResGetSignup, *response.Error) {
	return api.ResGetSignup{Url: signupUrl}, nil
}

func PostSignupCallback(ctx echo.Context, db db.DBClient) (api.ResPostSignupCallback, *response.Error) {
	var req api.ReqPostSignupCallback
	ctx.Bind(&req)
	studentNumber, err := getStudentNumberfromGoogle(req.Code)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	err = createUser(studentNumber, db)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	userID, derr := users.IDFromStudentNumber(studentNumber, db)
	if derr != nil {
		return api.ResPostSignupCallback{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザーIDの取得に失敗しました", Log: derr.Error()}
	}
	jwt, derr := validator.CreateToken(userID)
	if derr != nil {
		return api.ResPostSignupCallback{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "トークンの生成に失敗しました", Log: derr.Error()}
	}
	return api.ResPostSignupCallback{Jwt: jwt}, nil
}

func getStudentNumberfromGoogle(code string) (string, *response.Error) {
	token, err := gcpConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", token.AccessToken), nil)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	userInfo := UserInfoResponse{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	if err := userInfo.validate(); err != nil {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "使用出来ないアカウントです", Log: err.Error()}
	}
	studentNumber := strings.TrimSuffix(userInfo.Email, emailSuffix)
	return studentNumber, nil
}

func createUser(studentNumber string, db db.DBClient) *response.Error {
	query, err := db.Query.ReadFile("sql/users/create_users.sql")
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	_, err = db.Client.Exec(string(query), &params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザー登録済みです", Log: err.Error()}
			}
		}
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return nil
}
