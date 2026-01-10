package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type workWithRelationsRow struct {
	WorkId         string  `db:"work_id"`
	WorkName       string  `db:"work_name"`
	WorkUpdatedAt  string  `db:"work_updated_at"`
	AuthorUserId   *string `db:"author_user_id"`
	AuthorUsername *string `db:"author_username"`
	AuthorIconUrl  *string `db:"author_icon_url"`
	TagId          *string `db:"tag_id"`
	TagName        *string `db:"tag_name"`
}

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
	if res.Works == nil {
		res.Works = []api.ResGetWorkWorkObjectWork{}
	}
	return res, nil
}

type workOverview struct {
	Authors []workObjectAuthor
	Name    string `db:"name"`
	Tags    []workObjectTag
	WorkId  string `db:"work_id"`
}

func getWorkList(dbClient db.Client, offset *int, authorId *string) ([]workOverview, *response.Error) {
	params := struct {
		Offset   *int    `twowaysql:"offset"`
		AuthorId *string `twowaysql:"authorId"`
	}{
		Offset:   offset,
		AuthorId: authorId,
	}
	rows := []workWithRelationsRow{}
	err := dbClient.Select(&rows, "sql/work/select_work_with_relations.sql", &params)
	if err != nil {
		return []workOverview{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return aggregateWorkListRows(rows), nil
}

func aggregateWorkListRows(rows []workWithRelationsRow) []workOverview {
	workMap := make(map[string]*workOverview)
	workOrder := make([]string, 0)
	authorMap := make(map[string]map[string]bool)
	tagMap := make(map[string]map[string]bool)

	for _, row := range rows {
		if _, exists := workMap[row.WorkId]; !exists {
			workMap[row.WorkId] = &workOverview{
				WorkId:  row.WorkId,
				Name:    row.WorkName,
				Authors: []workObjectAuthor{},
				Tags:    []workObjectTag{},
			}
			workOrder = append(workOrder, row.WorkId)
			authorMap[row.WorkId] = make(map[string]bool)
			tagMap[row.WorkId] = make(map[string]bool)
		}

		// 作者情報を追加
		if row.AuthorUserId != nil && !authorMap[row.WorkId][*row.AuthorUserId] {
			author := workObjectAuthor{
				UserId: *row.AuthorUserId,
			}
			if row.AuthorUsername != nil {
				author.Username = *row.AuthorUsername
			}
			if row.AuthorIconUrl != nil {
				author.IconUrl = *row.AuthorIconUrl
			}
			workMap[row.WorkId].Authors = append(workMap[row.WorkId].Authors, author)
			authorMap[row.WorkId][*row.AuthorUserId] = true
		}

		// タグ情報を追加
		if row.TagId != nil && !tagMap[row.WorkId][*row.TagId] {
			tag := workObjectTag{
				TagId: *row.TagId,
			}
			if row.TagName != nil {
				tag.Name = *row.TagName
			}
			workMap[row.WorkId].Tags = append(workMap[row.WorkId].Tags, tag)
			tagMap[row.WorkId][*row.TagId] = true
		}
	}

	result := make([]workOverview, 0, len(workMap))
	for _, workId := range workOrder {
		result = append(result, *workMap[workId])
	}

	return result
}
