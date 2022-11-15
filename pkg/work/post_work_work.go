package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PostWorkWork(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostWorkWork) (api.ResGetWorkWorkWorkId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	requestBody.Auther = utils.GetUniqueString(append(requestBody.Auther, userId))
	requestBody.File = utils.GetUniqueString(requestBody.File)
	requestBody.Tag = utils.GetUniqueString(requestBody.Tag)
	workId, err := createWork(dbClient, requestBody)
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

func createWork(dbClient db.TransactionClient, requestBody api.ReqPostWorkWork) (string, *response.Error) {
	params := struct {
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	_, rerr := dbClient.Exec("sql/work/insert_work.sql", &params, true)
	if rerr != nil {
		if mysqlErr, ok := rerr.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "既に登録済みのタグです", Log: rerr.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	workId, rerr := dbClient.GetId()
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return workId, nil
}

func createWorkAuther(dbClient db.TransactionClient, workId string, autherIds []string) *response.Error {
	for _, autherId := range autherIds {
		params := struct {
			WorkId string `twowaysql:"workId"`
			UserId string `twowaysql:"userId"`
		}{
			WorkId: workId,
			UserId: autherId,
		}
		_, rerr := dbClient.Exec("sql/work/insert_work_user.sql", &params, false)
		if rerr != nil {
			return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
		}
	}
	return nil
}

func createWorkFile(dbClient db.TransactionClient, workId string, fileIds []string) *response.Error {
	for _, fileId := range fileIds {
		params := struct {
			WorkId string `twowaysql:"workId"`
			FileId string `twowaysql:"fileId"`
		}{
			WorkId: workId,
			FileId: fileId,
		}
		_, rerr := dbClient.Exec("sql/work/insert_work_file.sql", &params, false)
		if rerr != nil {
			return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
		}
	}
	return nil
}

func createWorkWorkTag(dbClient db.TransactionClient, workId string, tagIds []string) *response.Error {
	for _, tagId := range tagIds {
		params := struct {
			WorkId string `twowaysql:"workId"`
			TagId  string `twowaysql:"tagId"`
		}{
			WorkId: workId,
			TagId:  tagId,
		}
		_, rerr := dbClient.Exec("sql/work/insert_work_work_tag.sql", &params, false)
		if rerr != nil {
			return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
		}
	}
	return nil
}
