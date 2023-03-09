package blog

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutBlogBlogBlogId(ctx echo.Context, dbClient db.TransactionClient, blogId string, requestBody api.ReqPutBlogBlogBlogId) (api.ResGetBlogBlogBlogId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	permission, err := checkBlogAuthor(dbClient, blogId, userId)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	if !permission {
		return api.ResGetBlogBlogBlogId{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "編集権限がありません", Log: "no edit permission"}
	}
	err = deleteBlogBlogTag(dbClient, blogId)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	err = updateBlog(dbClient, blogId, requestBody)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	err = createBlogBlogTag(dbClient, blogId, requestBody.Tags)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	return GetBlogBlogBlogId(ctx, dbClient, blogId)
}

func checkBlogAuthor(dbClient db.Client, blogId string, userId string) (bool, *response.Error) {
	author, err := getBlogBlogAuthorId(dbClient, blogId)
	if err != nil {
		return false, err
	}
	if author.UserId == userId {
		return true, nil
	}
	return false, nil
}

func updateBlog(dbClient db.TransactionClient, blogId string, requestBody api.ReqPutBlogBlogBlogId) *response.Error {
	params := struct {
		BlogId   string `twowaysql:"blogId"`
		Title    string `twowaysql:"title"`
		Content  string `twowaysql:"content"`
		IsPublic bool   `twowaysql:"is_public"`
	}{
		BlogId:   blogId,
		Title:    requestBody.Title,
		Content:  requestBody.Content,
		IsPublic: requestBody.IsPublic,
	}
	_, rerr := dbClient.Exec("sql/blog/update_blog.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
