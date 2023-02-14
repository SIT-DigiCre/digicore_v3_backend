package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetGroupGroupId(ctx echo.Context, dbClient db.Client, groupId string) (api.ResGetGroupGroupId, *response.Error) {
	res := api.ResGetGroupGroupId{}
	userId := ctx.Get("user_id").(string)
	groupDetail, err := getGroupFromGroupId(dbClient, groupId, userId)
	if err != nil {
		return api.ResGetGroupGroupId{}, err
	}
	rerr := copier.Copy(&res, &groupDetail)
	if rerr != nil {
		return api.ResGetGroupGroupId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループの取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type groupDetail struct {
	GroupId     string `db:"group_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	UserCount   int    `db:"user_count"`
	Joinable    bool   `db:"joinable"`
	Joined      bool   `db:"joined"`
	Users       []groupDetailObjectUser
}

type groupDetailObjectUser struct {
	Name     string `db:"username"`
	UserIcon string `db:"icon_url"`
	UserId   string `db:"user_id"`
}

func getGroupFromGroupId(dbClient db.Client, groupId string, userId string) (groupDetail, *response.Error) {
	params := struct {
		GroupId string `twowaysql:"groupId"`
		UserId  string `twowaysql:"userId"`
	}{
		GroupId: groupId,
		UserId:  userId,
	}
	eventDetails := []groupDetail{}
	err := dbClient.Select(&eventDetails, "sql/group/select_group_from_group_id.sql", &params)
	if err != nil {
		return groupDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループの取得に失敗しました", Log: err.Error()}
	}
	if len(eventDetails) == 0 {
		return groupDetail{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "グループの取得に失敗しました", Log: "no rows in result"}
	}
	groupUsers := []groupDetailObjectUser{}
	err = dbClient.Select(&groupUsers, "sql/group/select_group_user_from_group_id.sql", &params)
	if err != nil {
		return groupDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループの取得に失敗しました", Log: err.Error()}
	}
	eventDetails[0].Users = groupUsers
	return eventDetails[0], nil
}
