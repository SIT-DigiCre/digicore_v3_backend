package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PutAdminChangeStudentNumber(_ echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutAdminChangeStudentNumber) (api.BlankSuccess, *response.Error) {
	params := struct {
		UserId        string `twowaysql:"userId"`
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		UserId:        requestBody.UserId,
		StudentNumber: requestBody.StudentNumber,
	}

	result, err := dbClient.Exec("sql/admin/update_user_student_number.sql", &params, false)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "指定された学籍番号は既に使用されています",
				Log:     err.Error(),
			}
		}
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     err.Error(),
		}
	}

	rowsAffected, rerr := result.RowsAffected()
	if rerr != nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     rerr.Error(),
		}
	}
	if rowsAffected == 0 {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "更新対象のユーザーが存在しないか、学籍番号が変更されていません",
			Log:     "user not found or student number unchanged",
		}
	}

	return api.BlankSuccess{Success: true}, nil
}
