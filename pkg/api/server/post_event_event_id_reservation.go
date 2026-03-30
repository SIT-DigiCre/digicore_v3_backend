package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/labstack/echo/v4"
)

func (s *server) PostEventEventIdReservation(ctx echo.Context, eventId string) error {
	// リクエストボディの解析
	var requestBody api.PostEventEventIdReservationJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "リクエストボディの解析に失敗しました", Log: err.Error()})
	}

	// バリデーション
	err := validator.Validate(requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	dbTranisactionClient, rerr := db.OpenTransaction()
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}
	defer func() {
		dbTranisactionClient.Rollback()
	}()

	// 管理者であるか確認
	userIdInterface := ctx.Get("user_id")
	if userIdInterface == nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 401, Level: "Info", Message: "認証が必要です", Log: "user_id not found in context"})
	}
	userId, ok := userIdInterface.(string)
	if !ok {
		return response.ErrorResponse(ctx, &response.Error{Code: 500, Level: "Error", Message: "ユーザーIDの型変換に失敗しました", Log: "user_id is not string"})
	}
	isAdmin, err := admin.CheckUserIsAdmin(&dbTranisactionClient, userId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	if !isAdmin {
		return response.ErrorResponse(ctx, &response.Error{Code: 403, Level: "Info", Message: "管理者権限が必要です", Log: "user is not admin"})
	}

	res, rerr := event.PostEventEventIdReservation(ctx, &dbTranisactionClient, eventId, requestBody)
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}

	if cerr := dbTranisactionClient.Commit(); cerr != nil {
		return response.ErrorResponse(ctx, cerr)
	}

	return response.SuccessResponse(ctx, res)
}
