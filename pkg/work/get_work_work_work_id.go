package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type workDetailWithRelationsRow struct {
	WorkId          string  `db:"work_id"`
	WorkName        string  `db:"work_name"`
	WorkDescription string  `db:"work_description"`
	AuthorUserId    *string `db:"author_user_id"`
	AuthorUsername  *string `db:"author_username"`
	AuthorIconUrl   *string `db:"author_icon_url"`
	TagId           *string `db:"tag_id"`
	TagName         *string `db:"tag_name"`
	FileId          *string `db:"file_id"`
	FileName        *string `db:"file_name"`
}

func GetWorkWorkWorkId(ctx echo.Context, dbClient db.Client, workId string) (api.ResGetWorkWorkWorkId, *response.Error) {
	res := api.ResGetWorkWorkWorkId{}
	work, err := getWorkFromTagId(dbClient, workId)
	if err != nil {
		return api.ResGetWorkWorkWorkId{}, err
	}
	rerr := copier.Copy(&res, &work)
	if rerr != nil {
		return api.ResGetWorkWorkWorkId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	if res.Authors == nil {
		res.Authors = []api.ResGetWorkWorkWorkIdObjectAuthor{}
	}
	if res.Files == nil {
		res.Files = []api.ResGetWorkWorkWorkIdObjectFile{}
	}
	if res.Tags == nil {
		res.Tags = []api.ResGetWorkWorkWorkIdObjectTag{}
	}
	return res, nil
}

type work struct {
	Authors     []workObjectAuthor `db:"author"`
	Description string             `db:"description"`
	Files       []workObjectFile   `db:"file"`
	Name        string             `db:"name"`
	Tags        []workObjectTag    `db:"tag"`
	WorkId      string             `db:"work_id"`
}

type workObjectAuthor struct {
	UserId   string `db:"user_id"`
	Username string `db:"username"`
	IconUrl  string `db:"icon_url"`
}

type workObjectTag struct {
	TagId string `db:"tag_id"`
	Name  string `db:"name"`
}

type workObjectFile struct {
	FileId string `db:"file_id"`
	Name   string `db:"name"`
}

func getWorkFromTagId(dbClient db.Client, workId string) (work, *response.Error) {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	rows := []workDetailWithRelationsRow{}
	rerr := dbClient.Select(&rows, "sql/work/select_work_from_work_id_with_relations.sql", &params)
	if rerr != nil {
		return work{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	if len(rows) == 0 {
		return work{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "作品が存在しません", Log: "no rows in result"}
	}
	return mapRowsToWorkDetail(rows), nil
}

func mapRowsToWorkDetail(rows []workDetailWithRelationsRow) work {
	firstRow := rows[0]
	result := work{
		WorkId:      firstRow.WorkId,
		Name:        firstRow.WorkName,
		Description: firstRow.WorkDescription,
		Authors:     []workObjectAuthor{},
		Tags:        []workObjectTag{},
		Files:       []workObjectFile{},
	}

	authorMap := make(map[string]bool)
	tagMap := make(map[string]bool)
	fileMap := make(map[string]bool)

	for _, row := range rows {
		// 作者情報を追加
		if row.AuthorUserId != nil && !authorMap[*row.AuthorUserId] {
			author := workObjectAuthor{
				UserId: *row.AuthorUserId,
			}
			if row.AuthorUsername != nil {
				author.Username = *row.AuthorUsername
			}
			if row.AuthorIconUrl != nil {
				author.IconUrl = *row.AuthorIconUrl
			}
			result.Authors = append(result.Authors, author)
			authorMap[*row.AuthorUserId] = true
		}

		// タグ情報を追加
		if row.TagId != nil && !tagMap[*row.TagId] {
			tag := workObjectTag{
				TagId: *row.TagId,
			}
			if row.TagName != nil {
				tag.Name = *row.TagName
			}
			result.Tags = append(result.Tags, tag)
			tagMap[*row.TagId] = true
		}

		// ファイル情報を追加
		if row.FileId != nil && !fileMap[*row.FileId] {
			file := workObjectFile{
				FileId: *row.FileId,
			}
			if row.FileName != nil {
				file.Name = *row.FileName
			}
			result.Files = append(result.Files, file)
			fileMap[*row.FileId] = true
		}
	}

	return result
}

func getWorkWorkAuthorList(dbClient db.Client, workId string) ([]workObjectAuthor, *response.Error) {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	workAuthors := []workObjectAuthor{}
	err := dbClient.Select(&workAuthors, "sql/work/select_work_work_author.sql", &params)
	if err != nil {
		return []workObjectAuthor{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return workAuthors, nil
}
