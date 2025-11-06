package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostGroup(ctx echo.Context, dbClient db.TransactionClient, req api.ReqPostGroup) (api.ResPostGroup, *response.Error) {
	userId := ctx.Get("user_id").(string)

	// isAdminGroup=trueの場合、リクエストユーザーが管理者かどうかを確認
	if req.IsAdminGroup {
		isAdmin, err := checkUserIsAdmin(dbClient, userId)
		if err != nil {
			return api.ResPostGroup{}, err
		}
		if !isAdmin {
			return api.ResPostGroup{}, &response.Error{
				Code:    http.StatusForbidden,
				Level:   "Info",
				Message: "管理者グループを作成する権限がありません",
				Log:     "user is not admin",
			}
		}
	}

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

	// isAdminGroup=trueの場合、group_claimsテーブルにclaim="admin"を追加
	if req.IsAdminGroup {
		err = insertGroupClaim(dbClient, groupId, "admin")
		if err != nil {
			return api.ResPostGroup{}, err
		}
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

