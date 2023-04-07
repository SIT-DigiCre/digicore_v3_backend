package google_auth

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/authenticator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PostSignupCallback(ctx echo.Context, dbTransactionClient db.TransactionClient, requestBody api.ReqPostSignupCallback) (api.ResPostSignupCallback, *response.Error) {
	studentNumber, err := getStudentNumberfromGoogle(requestBody.Code, signupRedirectUrl)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	userId, err := createUser(studentNumber, dbTransactionClient)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	jwt, err := authenticator.CreateToken(userId)
	if err != nil {
		return api.ResPostSignupCallback{}, err
	}
	utils.NoticeMattermost(fmt.Sprintf("%s(%s)がデジコアに登録しました", studentNumber, userId))
	return api.ResPostSignupCallback{Jwt: jwt}, nil
}

func createUser(studentNumber string, dbClient db.TransactionClient) (string, *response.Error) {
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	_, err := dbClient.Exec("sql/user/insert_user.sql", &params, true)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザー登録済みです", Log: err.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	id, err := dbClient.GetId()
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	rerr := createDefaultUser(dbClient, id, studentNumber)
	if rerr != nil {
		return "", rerr
	}
	return id, nil
}

func createDefaultUser(dbClient db.TransactionClient, userId string, studentNumber string) *response.Error {
	enterYear, err := strconv.Atoi(studentNumber[2:4])
	if err != nil {
		enterYear = utils.GetSchoolYear()
	}
	schoolGrade := utils.GetSchoolYear() - 2000 - enterYear + 1
	if studentNumber[0] == 'm' {
		schoolGrade += 4
	} else if studentNumber[0] == 'n' {
		schoolGrade += 6
	}
	rerr := user.UpdateUserProfile(dbClient, userId, api.ReqPutUserMe{Username: studentNumber, SchoolGrade: schoolGrade, IconUrl: env.DefaultIconUrl})
	if rerr != nil {
		return rerr
	}
	return nil
}
