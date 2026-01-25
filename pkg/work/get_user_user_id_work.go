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

func getWorkListFromAuthorId(dbClient db.Client, userId string) ([]workOverview, *response.Error) {
	params := struct {
		AuthorId *string `twowaysql:"authorId"`
	}{
		AuthorId: &userId,
	}

	rows := []workWithRelationsRow{}
	err := dbClient.Select(&rows, "sql/work/select_work_with_relations.sql", &params)
	if err != nil {
		return []workOverview{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     err.Error(),
		}
	}

	return mapRowsToWorkList(rows), nil
}
