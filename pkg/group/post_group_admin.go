package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostGroupAdmin(ctx echo.Context, dbClient db.TransactionClient, req api.ReqPostGroupAdmin) (api.ResPostGroup, *response.Error) {
	userId := ctx.Get("user_id").(string)

	// グループを作成
	err := insertGroup(dbClient, req.Name, req.Description, req.Joinable)
	if err != nil {
		return api.ResPostGroup{}, err
	}

	// 生成されたIDを取得
	groupId, rerr := dbClient.GetId()
	if rerr != nil {
		return api.ResPostGroup{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     rerr.Error(),
		}
	}

	// 作成者をグループメンバーとして追加
	err = insertGroupUser(dbClient, userId, groupId)
	if err != nil {
		return api.ResPostGroup{}, err
	}

	// group_claimsテーブルにclaimを追加。infra権限を持つユーザーはあらゆるclaimのグループを作成可能
	err = insertGroupClaim(dbClient, groupId, req.Claim)
	if err != nil {
		return api.ResPostGroup{}, err
	}

	// レスポンスを返す
	res := api.ResPostGroup{
		GroupId:     groupId,
		Name:        req.Name,
		Description: req.Description,
		Joinable:    req.Joinable,
		UserCount:   1,
	}

	return res, nil
}
