package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetWorkTag(ctx echo.Context, dbClient db.Client, params api.GetWorkTagParams) (api.ResGetWorkTag, *response.Error) {
	res := api.ResGetWorkTag{}
	tag, err := getWorkTagList(dbClient, params.Offset)
	if err != nil {
		return api.ResGetWorkTag{}, err
	}
	rerr := copier.Copy(&res.Tags, &tag)
	if rerr != nil {
		return api.ResGetWorkTag{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type tagOverview struct {
	TagId string `db:"tag_id"`
	Name  string `db:"name"`
}

func getWorkTagList(dbClient db.Client, offset *int) ([]tagOverview, *response.Error) {
	params := struct {
		Offset *int `twowaysql:"offset"`
	}{
		Offset: offset,
	}
	tagOverviews := []tagOverview{}
	err := dbClient.Select(&tagOverviews, "sql/work/select_work_tag.sql", &params)
	if err != nil {
		return []tagOverview{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(tagOverviews) == 0 {
		return []tagOverview{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "タグが存在しません", Log: "no rows in result"}
	}
	return tagOverviews, nil
}
