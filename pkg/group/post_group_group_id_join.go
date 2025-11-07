package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostGroupGroupIdJoin(ctx echo.Context, dbClient db.TransactionClient, groupId string) (api.ResPostGroupGroupIdJoin, *response.Error) {
	userId := ctx.Get("user_id").(string)

	// グループが存在するか確認
	groupExists, err := checkGroupExists(dbClient, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdJoin{}, err
	}
	if !groupExists {
		return api.ResPostGroupGroupIdJoin{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "指定されたグループが存在しません",
			Log:     "group not found",
		}
	}

	// グループがadminグループかどうか確認
	isAdminGroup, err := checkGroupIsAdminGroup(dbClient, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdJoin{}, err
	}
	if isAdminGroup {
		return api.ResPostGroupGroupIdJoin{}, &response.Error{
			Code:    http.StatusForbidden,
			Level:   "Info",
			Message: "このグループには参加できません",
			Log:     "cannot join admin group",
		}
	}

	// グループのjoinableがtrueかどうか確認
	joinable, err := checkGroupJoinable(dbClient, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdJoin{}, err
	}
	if !joinable {
		return api.ResPostGroupGroupIdJoin{}, &response.Error{
			Code:    http.StatusForbidden,
			Level:   "Info",
			Message: "このグループには参加できません",
			Log:     "group is not joinable",
		}
	}

	// リクエストユーザーが既にグループに所属していないか確認
	alreadyMember, err := checkUserIsGroupMember(dbClient, userId, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdJoin{}, err
	}
	if alreadyMember {
		return api.ResPostGroupGroupIdJoin{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "既にグループに参加しています",
			Log:     "user already member of the group",
		}
	}

	// groups_usersテーブルに追加
	err = insertGroupUser(dbClient, userId, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdJoin{}, err
	}

	// groupsテーブルのuser_countをインクリメント
	err = incrementGroupUserCount(dbClient, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdJoin{}, err
	}

	// レスポンスを返す
	res := api.ResPostGroupGroupIdJoin{
		Message: "グループに参加しました",
	}

	return res, nil
}
