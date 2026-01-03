package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetActivityPlacePlaceCurrent(ctx echo.Context, dbClient db.Client, place string) (api.ResGetActivityPlacePlaceCurrent, *response.Error) {
	res := api.ResGetActivityPlacePlaceCurrent{}
	users, err := selectCurrentUsers(dbClient, place)
	if err != nil {
		return api.ResGetActivityPlacePlaceCurrent{}, err
	}
	rerr := copier.Copy(&res.Users, &users)
	if rerr != nil {
		return api.ResGetActivityPlacePlaceCurrent{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "現在在室中のユーザー一覧の取得に失敗しました",
			Log:     rerr.Error(),
		}
	}
	if res.Users == nil {
		res.Users = []api.ResGetActivityPlacePlaceCurrentObjectUser{}
	}
	return res, nil
}

type currentUser struct {
	UserId            string    `db:"user_id"`
	Username          string    `db:"username"`
	ShortIntroduction string    `db:"short_introduction"`
	IconUrl           string    `db:"icon_url"`
	CheckedInAt       time.Time `db:"checked_in_at"`
}

func selectCurrentUsers(dbClient db.Client, place string) ([]currentUser, *response.Error) {
	params := struct {
		Place string `twowaysql:"place"`
	}{
		Place: place,
	}
	users := []currentUser{}
	err := dbClient.Select(&users, "sql/activity/select_current_users.sql", &params)
	if err != nil {
		return []currentUser{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "現在在室中のユーザー一覧の取得に失敗しました",
			Log:     err.Error(),
		}
	}
	return users, nil
}
