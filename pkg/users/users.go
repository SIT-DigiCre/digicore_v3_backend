package users

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func IDFromStudentNumber(studentNumber string, db db.DBClient) (string, *response.Error) {
	query, err := db.Query.ReadFile("sql/users/select_id_from_student_number.sql")
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	users := []struct {
		ID string `db:"id"`
	}{}
	err = db.Client.Select(&users, string(query), &params)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	if len(users) != 1 {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: "一意制約に違反しているデータがあります"}
	}
	return users[0].ID, nil
}
