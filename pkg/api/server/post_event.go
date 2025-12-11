package server

import (
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

func (s *server) PostEvent(ctx echo.Context) error {
	// Bind into a raw struct with string dates to avoid validator panic
	var raw struct {
		Name              string `json:"name" validate:"required"`
		Description       string `json:"description" validate:"required"`
		StartDate         string `json:"start_date" validate:"required"`
		FinishDate        string `json:"finish_date" validate:"required"`
		ReservationStart  string `json:"reservation_start" validate:"required"`
		ReservationFinish string `json:"reservation_finish" validate:"required"`
		Capacity          int    `json:"capacity" validate:"required,min=1"`
	}

	if err := ctx.Bind(&raw); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "リクエストボディの解析に失敗しました。正しい形式で送信してください", Log: err.Error()})
	}

	// Simple manual validation to avoid validator tag parse issues
	if raw.Name == "" || raw.Description == "" {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "必須フィールドが不足しています", Log: "missing required field"})
	}
	if raw.Capacity <= 0 {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "capacity は 1 以上である必要があります", Log: "invalid capacity"})
	}

	// parse datetime fields
	sd, err := time.Parse(time.RFC3339, raw.StartDate)
	if err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "start_date の形式が不正です", Log: err.Error()})
	}
	fd, err := time.Parse(time.RFC3339, raw.FinishDate)
	if err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "finish_date の形式が不正です", Log: err.Error()})
	}
	rs, err := time.Parse(time.RFC3339, raw.ReservationStart)
	if err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "reservation_start の形式が不正です", Log: err.Error()})
	}
	rf, err := time.Parse(time.RFC3339, raw.ReservationFinish)
	if err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "reservation_finish の形式が不正です", Log: err.Error()})
	}

	requestBody := api.PostEventJSONBody{
		Capacity:          raw.Capacity,
		Description:       raw.Description,
		FinishDate:        fd,
		Name:              raw.Name,
		ReservationFinish: rf,
		ReservationStart:  rs,
		StartDate:         sd,
	}

	dbTranisactionClient, rerr := db.OpenTransaction()
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}
	defer func() {
		_ = dbTranisactionClient.Rollback() // best-effort rollback, ignore error
	}()

	// Require authenticated admin user to create events
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

	/*if err := validator.Validate(&raw); err != nil {
		return response.ErrorResponse(ctx, err)
	}*/

	res, rerr := event.PostEvent(ctx, &dbTranisactionClient, requestBody)
	if rerr != nil {
		return response.ErrorResponse(ctx, rerr)
	}

	if cerr := dbTranisactionClient.Commit(); cerr != nil {
		return response.ErrorResponse(ctx, cerr)
	}

	return response.SuccessResponse(ctx, res)
}
