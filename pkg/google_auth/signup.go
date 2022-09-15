package google_auth

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func GetSignup(ctx echo.Context) (api.ResGetSignup, *response.Error) {
	return api.ResGetSignup{Url: signupUrl}, nil
}

func PostSignupCallback(ctx echo.Context, db db.DBClient) (api.ResPostSignupCallback, *response.Error) {

	var req api.ReqPostSignupCallback
	ctx.Bind(&req)
	studentNumber, err := getStudentIDfromGoogle(req.Code)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	id, err := createUser(studentNumber, db)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	return api.ResPostSignupCallback{Jwt: id}, nil
}

func getStudentIDfromGoogle(code string) (string, *response.Error) {
	return "", nil
}

func createUser(studentNumber string, db db.DBClient) (string, *response.Error) {
	query, err := db.Query.ReadFile("sql/users/create_users.sql")
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラー", Log: err.Error()}
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
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザー登録済みです", Log: err.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラー", Log: err.Error()}
	}
	return "", nil
}
