package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostGroupGroupIdUser(ctx echo.Context, dbClient db.TransactionClient, groupId string, req api.ReqPostGroupGroupIdUser) (api.ResPostGroupGroupIdUser, *response.Error) {
	requestUserId := ctx.Get("user_id").(string)

	// グループが存在するか確認
	groupExists, err := checkGroupExists(dbClient, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdUser{}, err
	}
	if !groupExists {
		return api.ResPostGroupGroupIdUser{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "指定されたグループが存在しません",
			Log:     "group not found",
		}
	}

	// リクエストユーザーがグループに所属しているか、またはinfra claimを持つか確認
	isMember, err := checkUserIsGroupMember(dbClient, requestUserId, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdUser{}, err
	}
	if !isMember {
		hasInfra, err := admin.CheckUserHasClaim(dbClient, requestUserId, "infra")
		if err != nil {
			return api.ResPostGroupGroupIdUser{}, err
		}
		if !hasInfra {
			return api.ResPostGroupGroupIdUser{}, &response.Error{
				Code:    http.StatusForbidden,
				Level:   "Info",
				Message: "グループに所属していないため、ユーザーを追加できません",
				Log:     "user is not member of the group",
			}
		}
	}

	// 追加対象ユーザーが存在するか確認
	userExists, err := checkUserExists(dbClient, req.UserId)
	if err != nil {
		return api.ResPostGroupGroupIdUser{}, err
	}
	if !userExists {
		return api.ResPostGroupGroupIdUser{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "指定されたユーザーが存在しません",
			Log:     "user not found",
		}
	}

	// 追加対象ユーザーが既にグループに所属していないか確認
	alreadyMember, err := checkUserIsGroupMember(dbClient, req.UserId, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdUser{}, err
	}
	if alreadyMember {
		return api.ResPostGroupGroupIdUser{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "指定されたユーザーは既にグループに所属しています",
			Log:     "user already member of the group",
		}
	}

	// groups_usersテーブルに追加
	err = insertGroupUser(dbClient, req.UserId, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdUser{}, err
	}

	// groupsテーブルのuser_countをインクリメント
	err = incrementGroupUserCount(dbClient, groupId)
	if err != nil {
		return api.ResPostGroupGroupIdUser{}, err
	}

	// レスポンスを返す
	res := api.ResPostGroupGroupIdUser{
		Message: "ユーザーをグループに追加しました",
	}

	return res, nil
}
