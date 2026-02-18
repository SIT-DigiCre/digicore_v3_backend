package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetGroup(ctx echo.Context, dbClient db.Client, params api.GetGroupParams) (api.ResGetGroup, *response.Error) {
	res := api.ResGetGroup{}
	userId := ctx.Get("user_id").(string)
	groups, err := getGroupList(dbClient, userId, params.Offset)
	if err != nil {
		return api.ResGetGroup{}, err
	}
	rerr := copier.Copy(&res.Groups, &groups)
	if rerr != nil {
		return api.ResGetGroup{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループ一覧の取得に失敗しました", Log: rerr.Error()}
	}
	if res.Groups == nil {
		res.Groups = []api.ResGetGroupObjectGroup{}
	}
	return res, nil
}

type group struct {
	GroupId      string `db:"group_id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	UserCount    int    `db:"user_count"`
	Joinable     bool   `db:"joinable"`
	Joined       bool   `db:"joined"`
	IsAdminGroup bool   `db:"is_admin_group"`
}

func getGroupList(dbClient db.Client, userId string, offset *int) ([]group, *response.Error) {
	params := struct {
		UserId      string   `twowaysql:"userId"`
		Offset      *int     `twowaysql:"offset"`
		AdminClaims []string `twowaysql:"adminClaims"`
	}{
		UserId:      userId,
		Offset:      offset,
		AdminClaims: admin.AdminClaims,
	}
	groups := []group{}
	err := dbClient.Select(&groups, "sql/group/select_group.sql", &params)
	if err != nil {
		return []group{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "グループ一覧の取得に失敗しました", Log: err.Error()}
	}
	return groups, nil
}
