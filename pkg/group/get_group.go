package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetGroup(ctx echo.Context, dbClient db.Client, params api.GetGroupParams) (api.ResGetGroup, *response.Error) {
	res := api.ResGetGroup{}
	userId := ctx.Get("user_id").(string)
	events, err := getGroupList(dbClient, userId, params.Offset)
	if err != nil {
		return api.ResGetGroup{}, err
	}
	rerr := copier.Copy(&res.Groups, &events)
	if rerr != nil {
		return api.ResGetGroup{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループ一覧の取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type group struct {
	GroupId     string `db:"group_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	UserCount   int    `db:"user_count"`
	Joinable    bool   `db:"joinable"`
	Joined      bool   `db:"joined"`
}

func getGroupList(dbClient db.Client, userId string, offset *int) ([]group, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
		Offset *int   `twowaysql:"offset"`
	}{
		UserId: userId,
		Offset: offset,
	}
	groups := []group{}
	err := dbClient.Select(&groups, "sql/group/select_group.sql", &params)
	if err != nil {
		return []group{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループ一覧の取得に失敗しました", Log: err.Error()}
	}
	return groups, nil
}
