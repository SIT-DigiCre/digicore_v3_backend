package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/blog"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) PostBlogTag(ctx echo.Context) error {
	var requestBody api.ReqPostBlogTag
	ctx.Bind(&requestBody)
	err := validator.Validate(requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	dbTranisactionClient, err := db.OpenTransaction()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}
	defer dbTranisactionClient.Rollback()

	res, err := blog.PostBlogTag(ctx, &dbTranisactionClient, requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	err = dbTranisactionClient.Commit()
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
