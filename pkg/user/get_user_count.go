package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserCount(ctx echo.Context, dbClient db.Client) (api.ResGetUserCount, *response.Error) {
	rows := []struct {
		Count int `db:"count"`
	}{}
	params := struct{}{}
	err := dbClient.Select(&rows, "sql/user/select_user_count.sql", &params)
	if err != nil {
		return api.ResGetUserCount{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(rows) == 0 {
		return api.ResGetUserCount{Count: 0}, nil
	}
	return api.ResGetUserCount{Count: rows[0].Count}, nil
}
