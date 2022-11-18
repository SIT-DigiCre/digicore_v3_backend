package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteWorkTagTagId(ctx echo.Context, dbClient db.TransactionClient, tagId string) (api.Success, *response.Error) {
	err := deleteWorkTagFromTagId(dbClient, tagId)
	if err != nil {
		return api.Success{}, err
	}

	return api.Success{Success: true}, nil
}

func deleteWorkTagFromTagId(dbClient db.TransactionClient, tagId string) *response.Error {
	params := struct {
		TagId string `twowaysql:"tagId"`
	}{
		TagId: tagId,
	}
	_, err := dbClient.Exec("sql/work/delete_work_tag_from_tag_id.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
