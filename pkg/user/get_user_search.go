package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserSearch(ctx echo.Context, dbClient db.Client, params api.GetUserSearchParams) (api.ResGetUser, *response.Error) {
	res := api.ResGetUser{}
	if params.Query == "" {
		return api.ResGetUser{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "検索クエリが空です", Log: "Empty query parameter provided"}
	}
	user, err := searchUserList(dbClient, params.Query)
	if err != nil {
		return api.ResGetUser{}, err
	}
	rerr := copier.Copy(&res.Users, &user)
	if rerr != nil {
		return api.ResGetUser{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	if res.Users == nil {
		res.Users = []api.ResGetUserObjectUser{}
	}
	return res, nil
}

func searchUserList(dbClient db.Client, query string) ([]userOverview, *response.Error) {
	params := struct {
		Query string `twowaysql:"query"`
	}{
		Query: query,
	}
	userOverviews := []userOverview{}
	err := dbClient.Select(&userOverviews, "sql/user/select_user_search.sql", &params)
	if err != nil {
		return []userOverview{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return userOverviews, nil
}
