package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/labstack/echo/v4"
)

func (s *server) PutEventEventId(ctx echo.Context, eventId string) error {
	// リクエストボディの解析
	var requestBody api.PutEventEventIdJSONRequestBody
	if err := ctx.Bind(&requestBody); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "リクエストボディの解析に失敗しました", Log: err.Error()})
	}

	dbTranisactionClient, rerr := db.OpenTransaction()
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}
	defer func() {
		dbTranisactionClient.Rollback()
	}()

	// 管理者であるか確認
	userId := ctx.Get("user_id").(string)
	isAdmin, err := admin.CheckUserIsAdmin(&dbTranisactionClient, userId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	if !isAdmin {
		return response.ErrorResponse(ctx, &response.Error{Code: 403, Level: "Info", Message: "管理者権限が必要です", Log: "user is not admin"})
	}

	res, rerr := event.PutEvent(ctx, &dbTranisactionClient, eventId, requestBody)
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}

	if cerr := dbTranisactionClient.Commit(); cerr != nil {
		return response.ErrorResponse(ctx, cerr)
	}

	return response.SuccessResponse(ctx, res)
}
