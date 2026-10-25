package work

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/google/uuid"
)

func setWorkThumbnail(client db.TransactionClient, workID, userID string, thumbnailID *string) *response.Error {
	if thumbnailID == nil {
		return nil
	}
	if *thumbnailID != "" {
		parsed, err := uuid.Parse(*thumbnailID)
		if err != nil {
			return thumbnailError(http.StatusBadRequest, "サムネイルIDが不正です")
		}
		id := parsed.String()
		thumbnailID = &id
		params := struct {
			ThumbnailId string `twowaysql:"thumbnailId"`
			WorkId      string `twowaysql:"workId"`
			UserId      string `twowaysql:"userId"`
		}{id, workID, userID}
		rows := []struct {
			Allowed bool `db:"allowed"`
		}{}
		err = client.Select(&rows, "sql/work/select_work_thumbnail_permission.sql", &params)
		if err != nil {
			return &response.Error{Code: 500, Level: "Error", Message: "DBエラーが発生しました", Log: err.Error()}
		}
		if len(rows) == 0 || !rows[0].Allowed {
			return thumbnailError(http.StatusBadRequest, "指定したサムネイルは利用できません")
		}
	}
	params := struct {
		WorkId      string `twowaysql:"workId"`
		ThumbnailId string `twowaysql:"thumbnailId"`
	}{workID, *thumbnailID}
	_, err := client.Exec("sql/work/update_work_thumbnail.sql", &params, false)
	if err != nil {
		return &response.Error{Code: 500, Level: "Error", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return nil
}
