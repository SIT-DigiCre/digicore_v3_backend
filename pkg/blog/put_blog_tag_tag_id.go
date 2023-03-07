package blog

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutBlogTagTagId(ctx echo.Context, dbClient db.TransactionClient, tagId string, requestBody api.ReqPutBlogTagTagId) (api.ResGetBlogTagTagId, *response.Error) {
	err := updateBlogTag(dbClient, tagId, api.ReqPutBlogTagTagId(requestBody))
	if err != nil {
		return api.ResGetBlogTagTagId{}, err
	}
	return GetBlogTagTagId(ctx, dbClient, tagId)
}

func updateBlogTag(dbClient db.TransactionClient, tagId string, requestBody api.ReqPutBlogTagTagId) *response.Error {
	params := struct {
		TagId       string `twowaysql:"tagId"`
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		TagId:       tagId,
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	_, rerr := dbClient.Exec("sql/blog/update_blog_tag_from_tag_id.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
