package work

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PostWorkTag(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostWorkTag) (api.ResGetWorkTagTagId, *response.Error) {
	tagId, err := createWorkTag(dbClient, requestBody)
	if err != nil {
		return api.ResGetWorkTagTagId{}, err
	}
	fmt.Print(tagId)

	return GetWorkTagTagId(ctx, dbClient, tagId)
}

func createWorkTag(dbClient db.TransactionClient, requestBody api.ReqPostWorkTag) (string, *response.Error) {
	params := struct {
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	_, rerr := dbClient.Exec("sql/work/insert_work_tag.sql", &params, true)
	if rerr != nil {
		if mysqlErr, ok := rerr.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "既に登録済みのタグです", Log: rerr.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	tagId, rerr := dbClient.GetId()
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return tagId, nil
}
