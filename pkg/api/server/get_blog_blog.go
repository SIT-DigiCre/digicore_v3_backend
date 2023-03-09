package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/blog"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func (s *server) GetBlogBlog(ctx echo.Context, params api.GetBlogBlogParams) error {
	dbClient := db.Open()

	res, err := blog.GetBlogBlog(ctx, &dbClient, params)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
