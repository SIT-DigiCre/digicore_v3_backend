package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

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
	return res, nil
}

type work struct {
	Authors     []workObjectAuthor `db:"author"`
	Description string             `db:"description"`
	File        []workObjectFile   `db:"file"`
	Name        string             `db:"name"`
	Tag         []workObjectTag    `db:"tag"`
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
	works := []work{}
	rerr := dbClient.Select(&works, "sql/work/select_work_from_work_id.sql", &params)
	if rerr != nil {
		return work{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	if len(works) == 0 {
		return work{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "作品が存在しません", Log: "no rows in result"}
	}
	workAuthors, err := getWorkWorkAuthorList(dbClient, workId)
	if err != nil {
		return work{}, err
	}
	works[0].Authors = workAuthors
	workTags, err := getWorkWorkTagList(dbClient, workId)
	if err != nil {
		return work{}, err
	}
	works[0].Tag = workTags
	workFiles, err := getWorkWorkFileList(dbClient, workId)
	if err != nil {
		return work{}, err
	}
	works[0].File = workFiles
	return works[0], nil
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

func getWorkWorkTagList(dbClient db.Client, workId string) ([]workObjectTag, *response.Error) {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	workTags := []workObjectTag{}
	err := dbClient.Select(&workTags, "sql/work/select_work_work_tag.sql", &params)
	if err != nil {
		return []workObjectTag{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return workTags, nil
}

func getWorkWorkFileList(dbClient db.Client, workId string) ([]workObjectFile, *response.Error) {
	params := struct {
		WorkId string `twowaysql:"workId"`
	}{
		WorkId: workId,
	}
	workFiles := []workObjectFile{}
	err := dbClient.Select(&workFiles, "sql/work/select_work_work_file.sql", &params)
	if err != nil {
		return []workObjectFile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return workFiles, nil
}
