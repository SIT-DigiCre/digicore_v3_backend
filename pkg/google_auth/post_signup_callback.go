package google_auth

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/authenticator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PostSignupCallback(ctx echo.Context, dbTransactionClient db.TransactionClient, requestBody api.ReqPostSignupCallback) (api.ResPostSignupCallback, *response.Error) {
	studentNumber, err := getStudentNumberfromGoogle(requestBody.Code, signupRedirectUrl)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	userID, err := createUser(studentNumber, dbTransactionClient)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	jwt, err := authenticator.CreateToken(userID)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	return api.ResPostSignupCallback{Jwt: jwt}, nil
}

func createUser(studentNumber string, dbClient db.TransactionClient) (string, *response.Error) {
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	_, err := dbClient.Exec("sql/users/insert_users.sql", &params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザー登録済みです", Log: err.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return "", nil
}
