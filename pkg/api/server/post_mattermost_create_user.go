package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/mattermost"
	"github.com/labstack/echo/v4"
)

func (s *server) PostMattermostCreateUser(ctx echo.Context) error {
	var requestBody api.ReqPostMattermostCreateuser
	ctx.Bind(&requestBody)
	err := validator.Validate(requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	dbClient := db.Open()

	res, err := mattermost.PostMattermostCreateUser(ctx, &dbClient, requestBody)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
