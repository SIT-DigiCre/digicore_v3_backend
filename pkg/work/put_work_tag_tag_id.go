package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutWorkTagTagId(ctx echo.Context, dbClient db.TransactionClient, tagId string, requestBody api.ReqPutWorkTagTagId) (api.ResGetWorkTagTagId, *response.Error) {
	err := updateWorkTag(dbClient, tagId, api.ReqPostWorkTag(requestBody))
	if err != nil {
		return api.ResGetWorkTagTagId{}, err
	}
	return GetWorkTagTagId(ctx, dbClient, tagId)
}

func updateWorkTag(dbClient db.TransactionClient, tagId string, requestBody api.ReqPostWorkTag) *response.Error {
	params := struct {
		TagId       string `twowaysql:"tagId"`
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		TagId:       tagId,
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	_, rerr := dbClient.Exec("sql/work/update_work_tag_from_tag_id.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
