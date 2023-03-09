package blog

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteBlogBlogBlogId(ctx echo.Context, dbClient db.TransactionClient, blogId string) (api.BlankSuccess, *response.Error) {
	err := deleteBlogBlogTag(dbClient, blogId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	err = deleteBlog(dbClient, blogId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	return api.BlankSuccess{Success: true}, nil
}

func deleteBlog(dbClient db.TransactionClient, blogId string) *response.Error {
	params := struct {
		BlogId string `twowaysql:"blogId"`
	}{
		BlogId: blogId,
	}
	_, rerr := dbClient.Exec("sql/blog/delete_blog.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}

func deleteBlogBlogTag(dbClient db.TransactionClient, blogId string) *response.Error {
	params := struct {
		BlogId string `twowaysql:"blogId"`
	}{
		BlogId: blogId,
	}
	_, rerr := dbClient.Exec("sql/blog/delete_blog_blog_tag.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
