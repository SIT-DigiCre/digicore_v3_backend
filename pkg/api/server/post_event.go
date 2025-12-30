package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

func (s *server) PostEvent(ctx echo.Context) error {
	// リクエストボディの解析
	var requestBody api.PostEventJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "リクエストボディの解析に失敗しました", Log: err.Error()})
	}

	// 空データのチェック
	if requestBody.Name == "" || requestBody.Description == "" {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "必須フィールドが不足しています", Log: "missing required field"})
	}

	dbTranisactionClient, rerr := db.OpenTransaction()
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}
	defer func() {
		dbTranisactionClient.Rollback()
	}()

	// 管理者であるか確認
	var userId string
	if v := ctx.Get("user_id"); v != nil {
		if s, ok := v.(string); ok {
			userId = s
		}
	}
	if userId == "" {
		return response.ErrorResponse(ctx, &response.Error{Code: 401, Level: "Info", Message: "ログインが必要です", Log: "missing user_id"})
	}
	profile, perr := user.GetUserProfileFromUserId(&dbTranisactionClient, userId)
	if perr != nil {
		return response.ErrorResponse(ctx, perr)
	}
	isAdmin := profile.IsAdmin
	if !isAdmin {
		return response.ErrorResponse(ctx, &response.Error{Code: 403, Level: "Info", Message: "管理者権限が必要です", Log: "user is not admin"})
	}

	res, rerr := event.PostEvent(ctx, &dbTranisactionClient, requestBody)
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}

	if cerr := dbTranisactionClient.Commit(); cerr != nil {
		return response.ErrorResponse(ctx, cerr)
	}

	return response.SuccessResponse(ctx, res)
}
