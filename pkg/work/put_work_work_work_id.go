package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/util"
	"github.com/labstack/echo/v4"
)

func PutWorkWorkWorkId(ctx echo.Context, dbClient db.TransactionClient, workId string, requestBody api.ReqPutWorkWorkWorkId) (api.ResGetWorkWorkWorkId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	requestBody.Authors = util.GetUniqueString(append(requestBody.Authors, userId))
	permission, err := checkWorkAuthor(dbClient, workId, userId)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	if !permission {
		return api.ResGetWorkWorkWorkId{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "編集権限がありません", Log: "no edit permission"}
	}
	err = deleteWorkAuthor(dbClient, workId)
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
	err = createWorkAuthor(dbClient, workId, requestBody.Authors)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = createWorkFile(dbClient, workId, requestBody.Files)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	err = createWorkWorkTag(dbClient, workId, requestBody.Tags)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	return GetWorkWorkWorkId(ctx, dbClient, workId)
}

func checkWorkAuthor(dbClient db.Client, workId string, userId string) (bool, *response.Error) {
	authors, err := getWorkWorkAuthorList(dbClient, workId)
	if err != nil {
		return false, err
	}
	for _, author := range authors {
		if author.UserId == userId {
			return true, nil
		}
	}
	return false, nil
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
