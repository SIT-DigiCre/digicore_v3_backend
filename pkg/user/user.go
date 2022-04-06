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

func CreateDefault(db *sql.DB, id string, name string) error {
	enterYear, err := strconv.Atoi(name[2:4])
	if err != nil {
		return err
	}
	fmt.Printf("%d", enterYear)
	schoolGrade := time.Now().Year() - 2000 - enterYear + 1
	_, err = db.Exec(`INSERT INTO UserProfile (user_id, username, school_grade, icon_url, short_self_introduction, discord_userid) VALUES (UUID_TO_BIN(?), ?, ?, ?, 'デジクリ入りました', ?)`, id, name, schoolGrade, env.DefaultIconURL, "null")
	if err != nil {
		return err
	}
	return nil
}
