package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUser(ctx echo.Context, dbClient db.Client, params api.GetUserParams) (api.ResGetUser, *response.Error) {
	res := api.ResGetUser{}
	user, err := getUserList(dbClient, params.Offset, params.Seed)
	if err != nil {
		return api.ResGetUser{}, err
	}
	rerr := copier.Copy(&res.Users, &user)
	if rerr != nil {
		return api.ResGetUser{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type userOverview struct {
	IconUrl           string `db:"icon_url"`
	ShortIntroduction string `db:"short_introduction"`
	UserId            string `db:"user_id"`
	Username          string `db:"username"`
}

func getUserList(dbClient db.Client, offset *int, seed *int) ([]userOverview, *response.Error) {
	params := struct {
		Offset *int `twowaysql:"offset"`
		Seed   *int `twowaysql:"seed"`
	}{
		Offset: offset,
		Seed:   seed,
	}
	userOverviews := []userOverview{}
	err := dbClient.Select(&userOverviews, "sql/user/select_user_profile.sql", &params)
	if err != nil {
		return []userOverview{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return userOverviews, nil
}
