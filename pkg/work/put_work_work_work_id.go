package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutWorkWorkWorkId(ctx echo.Context, dbClient db.TransactionClient, workId string, requestBody api.ReqPutWorkWorkWorkId) (api.ResGetWorkWorkWorkId, *response.Error) {
	err := deleteWorkAuther(dbClient, workId)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = deleteWorkFile(dbClient, workId)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = deleteWorkWorkTag(dbClient, workId)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = updateWork(dbClient, workId, requestBody)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = createWorkAuther(dbClient, workId, requestBody.Auther)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = createWorkFile(dbClient, workId, requestBody.File)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = createWorkWorkTag(dbClient, workId, requestBody.Tag)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	return GetWorkWorkWorkId(ctx, dbClient, workId)
}

func updateWork(dbClient db.TransactionClient, workId string, requestBody api.ReqPutWorkWorkWorkId) *response.Error {
	params := struct {
		WorkId      string `twowaysql:"workId"`
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		WorkId:      workId,
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	_, rerr := dbClient.Exec("sql/work/update_work.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
