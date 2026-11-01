package server

import (
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/work"
	"github.com/labstack/echo/v4"
)

func (s *server) PostWorkThumbnail(ctx echo.Context) error {
	var body api.ReqPostWorkThumbnail
	if err := ctx.Bind(&body); err != nil {
		return response.ErrorResponse(ctx, &response.Error{Code: 400, Level: "Info", Message: "リクエストボディの解析に失敗しました", Log: err.Error()})
	}
	body.File = strings.NewReplacer("\r", "", "\n", "").Replace(body.File)
	if err := validator.Validate(body); err != nil {
		return response.ErrorResponse(ctx, err)
	}
	tx, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer tx.Rollback()
	res, err := work.PostWorkThumbnail(ctx, &tx, body)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	if err = tx.Commit(); err != nil {
		return response.ErrorResponse(ctx, err)
	}
	return response.SuccessResponse(ctx, res)
}
