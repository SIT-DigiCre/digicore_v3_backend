package users

import (
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func IDFromStudentNumber(studentNumber string, db db.DBClient) (string, error) {
	query, err := db.Query.ReadFile("sql/users/id_from_student_number.sql")
	if err != nil {
		return "", err
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
		return "", err
	}
	if len(users) != 1 {
		return "", fmt.Errorf("不明なエラー")
	}
	return users[0].ID, nil
}
