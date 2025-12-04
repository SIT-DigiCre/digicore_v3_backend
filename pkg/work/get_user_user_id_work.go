package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserUserIdWork(ctx echo.Context, dbClient db.Client, userId string) (api.ResGetWorkWork, *response.Error) {
	res := api.ResGetWorkWork{}

	workOverviews, err := getWorkListFromAuthorId(dbClient, userId)
	if err != nil {
		return api.ResGetWorkWork{}, err
	}

	rerr := copier.Copy(&res.Works, &workOverviews)
	if rerr != nil {
		return api.ResGetWorkWork{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     rerr.Error(),
		}
	}
	if res.Works == nil {
		res.Works = []api.ResGetWorkWorkObjectWork{}
	}

	return res, nil
}

// 指定されたユーザーが作者として含まれる作品一覧を取得する
func getWorkListFromAuthorId(dbClient db.Client, userId string) ([]workOverview, *response.Error) {
	params := struct {
		AuthorId *string `twowaysql:"authorId"`
	}{
		AuthorId: &userId,
	}

	workOverviews := []workOverview{}
	err := dbClient.Select(&workOverviews, "sql/work/select_work.sql", &params)
	if err != nil {
		return []workOverview{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     err.Error(),
		}
	}

	for i := range workOverviews {
		workId := workOverviews[i].WorkId
		workAuthors, err := getWorkWorkAuthorList(dbClient, workId)
		if err != nil {
			return []workOverview{}, err
		}
		workOverviews[i].Authors = workAuthors
		workTags, err := getWorkWorkTagList(dbClient, workId)
		if err != nil {
			return []workOverview{}, err
		}
		workOverviews[i].Tags = workTags
	}

	return workOverviews, nil
}
