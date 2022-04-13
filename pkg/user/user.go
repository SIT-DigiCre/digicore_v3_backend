package user

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Context struct {
	DB *sql.DB
}

type UserStudentNumber struct {
	StudentNumber string
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}

func GetUserId(e *echo.Context) (string, error) {
	user := (*e).Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["uuid"].(string)
	return id, nil
}

func GetStudentNumber(c Context, userId string) string {
	userStudentNumber := UserStudentNumber{}
	_ = c.DB.QueryRow("SELECT student_number FROM User WHERE id = UUID_TO_BIN(?)", userId).
		Scan(&userStudentNumber.StudentNumber)
	return userStudentNumber.StudentNumber
}

func CreateDefault(db *sql.DB, id string, name string) error {
	enterYear, err := strconv.Atoi(name[2:4])
	if err != nil {
		return err
	}
	fmt.Printf("%d", enterYear)
	schoolGrade := time.Now().Year() - 2000 - enterYear + 1
	_, err = db.Exec(`INSERT INTO UserProfile (user_id, username, school_grade, icon_url, active_limit) VALUES (UUID_TO_BIN(?), ?, ?, ?, CURRENT_DATE)`, id, name, schoolGrade, env.DefaultIconURL)
	if err != nil {
		return err
	}
	return nil
}
