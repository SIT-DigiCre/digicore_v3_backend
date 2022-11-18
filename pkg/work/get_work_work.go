package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetWorkWork(ctx echo.Context, dbClient db.Client, params api.GetWorkWorkParams) (api.ResGetWorkWork, *response.Error) {
	res := api.ResGetWorkWork{}
	work, err := getWorkList(dbClient, params.Offset, params.AuthorId)
	if err != nil {
		return api.ResGetWorkWork{}, err
	}
	rerr := copier.Copy(&res.Works, &work)
	if rerr != nil {
		return api.ResGetWorkWork{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type workOverview struct {
	Author []workObjectAuthor
	Name   string `db:"name"`
	Tag    []workObjectTag
	WorkId string `db:"work_id"`
}

func getWorkList(dbClient db.Client, offset *int, authorId *string) ([]workOverview, *response.Error) {
	params := struct {
		Offset   *int    `twowaysql:"offset"`
		AuthorId *string `twowaysql:"authorId"`
	}{
		Offset:   offset,
		AuthorId: authorId,
	}
	workOverviews := []workOverview{}
	err := dbClient.Select(&workOverviews, "sql/work/select_work.sql", &params)
	if err != nil {
		return []workOverview{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(workOverviews) == 0 {
		return []workOverview{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "作品が存在しません", Log: "no rows in result"}
	}
	for i := range workOverviews {
		workId := workOverviews[i].WorkId
		workAuthors, err := getWorkWorkAuthorList(dbClient, workId)
		if err != nil {
			return []workOverview{}, err
		}
		workOverviews[i].Author = workAuthors
		workTags, err := getWorkWorkTagList(dbClient, workId)
		if err != nil {
			return []workOverview{}, err
		}
		workOverviews[i].Tag = workTags
	}
	return workOverviews, nil
}
