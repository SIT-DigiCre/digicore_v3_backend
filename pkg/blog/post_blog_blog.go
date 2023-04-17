package blog

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PostBlogBlog(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostBlogBlog) (api.ResGetBlogBlogBlogId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	requestBody.Tags = utils.GetUniqueString(requestBody.Tags)
	blogId, err := createBlog(dbClient, requestBody, userId)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	err = createBlogBlogTag(dbClient, blogId, requestBody.Tags)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	return GetBlogBlogBlogId(ctx, dbClient, blogId)
}

func createBlog(dbClient db.TransactionClient, requestBody api.ReqPostBlogBlog, userId string) (string, *response.Error) {
	params := struct {
		Name     string `twowaysql:"name"`
		Content  string `twowaysql:"content"`
		IsPublic bool   `twowaysql:"isPublic"`
		UserId   string `twowaysql:"userId"`
		TopImage string `twowaysql:"topImage"`
	}{
		Name:     requestBody.Name,
		Content:  requestBody.Content,
		IsPublic: requestBody.IsPublic,
		UserId:   userId,
		TopImage: requestBody.TopImage,
	}
	_, rerr := dbClient.Exec("sql/blog/insert_blog.sql", &params, true)
	if rerr != nil {
		if mysqlErr, ok := rerr.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "既に登録済みのタグです", Log: rerr.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	blogId, rerr := dbClient.GetId()
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return blogId, nil
}

func createBlogBlogTag(dbClient db.TransactionClient, blogId string, tagIds []string) *response.Error {
	for _, tagId := range tagIds {
		params := struct {
			BlogId string `twowaysql:"blogId"`
			TagId  string `twowaysql:"tagId"`
		}{
			BlogId: blogId,
			TagId:  tagId,
		}
		_, rerr := dbClient.Exec("sql/blog/insert_blog_blog_tag.sql", &params, false)
		if rerr != nil {
			return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました1", Log: rerr.Error()}
		}
	}
	return nil
}
