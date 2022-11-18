package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteWorkWorkWorkId(ctx echo.Context, dbClient db.TransactionClient, workId string) (api.BlankSuccess, *response.Error) {
	err := deleteWorkAuthor(dbClient, workId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	err = deleteWorkFile(dbClient, workId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	err = deleteWorkWorkTag(dbClient, workId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	err = deleteWork(dbClient, workId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	return api.BlankSuccess{Success: true}, nil
}

func deleteWork(dbClient db.TransactionClient, workId string) *response.Error {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	_, rerr := dbClient.Exec("sql/work/delete_work.sql", &params, true)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}

func deleteWorkAuthor(dbClient db.TransactionClient, workId string) *response.Error {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	_, rerr := dbClient.Exec("sql/work/delete_work_user.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}

func deleteWorkFile(dbClient db.TransactionClient, workId string) *response.Error {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	_, rerr := dbClient.Exec("sql/work/delete_work_file.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}

	return nil
}

func deleteWorkWorkTag(dbClient db.TransactionClient, workId string) *response.Error {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	_, rerr := dbClient.Exec("sql/work/delete_work_work_tag.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
