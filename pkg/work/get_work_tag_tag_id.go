package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetWorkTagTagId(ctx echo.Context, dbClient db.Client, tagId string) (api.ResGetWorkTagTagId, *response.Error) {
	res := api.ResGetWorkTagTagId{}
	tag, err := getWorkTagFromTagId(dbClient, tagId)
	if err != nil {
		return api.ResGetWorkTagTagId{}, err
	}
	rerr := copier.Copy(&res, &tag)
	if rerr != nil {
		return api.ResGetWorkTagTagId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type tag struct {
	TagId       string `db:"tag_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func getWorkTagFromTagId(dbClient db.Client, tagId string) (tag, *response.Error) {
	params := struct {
		TagId string `twowaysql:"tagId"`
	}{
		TagId: tagId,
	}
	tags := []tag{}
	err := dbClient.Select(&tags, "sql/work/select_work_tag_from_tag_id.sql", &params)
	if err != nil {
		return tag{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(tags) == 0 {
		return tag{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "タグが存在しません", Log: "no rows in result"}
	}
	return tags[0], nil
}
