package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutWorkWorkWorkId(ctx echo.Context, dbClient db.TransactionClient, workId string, requestBody api.ReqPutWorkWorkWorkId) (api.ResGetWorkWorkWorkId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	requestBody.Auther = utils.GetUniqueString(append(requestBody.Auther, userId))
	permission, err := checkWorkAuther(dbClient, workId, userId)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	if !permission {
		return api.ResGetWorkWorkWorkId{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "編集権限がありません", Log: "no edit permission"}
	}
	err = deleteWorkAuther(dbClient, workId)
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

func checkWorkAuther(dbClient db.Client, workId string, userId string) (bool, *response.Error) {
	authers, err := getWorkWorkAutherList(dbClient, workId)
	if err != nil {
		return false, err
	}
	for _, auther := range authers {
		if auther.UserId == userId {
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
