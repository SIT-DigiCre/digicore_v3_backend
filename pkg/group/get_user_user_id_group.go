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

func GetUserUserIdGroup(ctx echo.Context, dbClient db.Client, userId string) (api.ResGetGroup, *response.Error) {
	res := api.ResGetGroup{}

	groups, err := getGroupListFromUserId(dbClient, userId)
	if err != nil {
		return api.ResGetGroup{}, err
	}

	rerr := copier.Copy(&res.Groups, &groups)
	if rerr != nil {
		return api.ResGetGroup{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "グループ一覧の取得に失敗しました",
			Log:     rerr.Error(),
		}
	}
	if res.Groups == nil {
		res.Groups = []api.ResGetGroupObjectGroup{}
	}

	return res, nil
}

// 指定されたユーザーが参加しているグループ一覧を取得する
func getGroupListFromUserId(dbClient db.Client, userId string) ([]group, *response.Error) {
	params := struct {
		UserId      string   `twowaysql:"userId"`
		AdminClaims []string `twowaysql:"adminClaims"`
	}{
		UserId:      userId,
		AdminClaims: admin.AdminClaims,
	}

	groups := []group{}
	err := dbClient.Select(&groups, "sql/group/select_group_from_user_id.sql", &params)
	if err != nil {
		return []group{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "グループ一覧の取得に失敗しました",
			Log:     err.Error(),
		}
	}

	return groups, nil
}
