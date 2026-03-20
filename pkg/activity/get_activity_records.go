package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetActivityRecords(_ echo.Context, dbClient db.Client, offset *int, limit *int) (api.ResGetActivityRecords, *response.Error) {
	res := api.ResGetActivityRecords{}

	defaultOffset := 0
	defaultLimit := 50
	if offset == nil {
		offset = &defaultOffset
	}
	if limit == nil {
		limit = &defaultLimit
	}

	records, err := selectActivityRecords(dbClient, *offset, *limit)
	if err != nil {
		return api.ResGetActivityRecords{}, err
	}

	total, err := selectActivityRecordsCount(dbClient)
	if err != nil {
		return api.ResGetActivityRecords{}, err
	}

	rerr := copier.Copy(&res.Records, &records)
	if rerr != nil {
		return api.ResGetActivityRecords{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "アクティビティレコード一覧のデータ変換に失敗しました",
			Log:     rerr.Error(),
		}
	}
	if res.Records == nil {
		res.Records = []api.ResGetActivityRecordsObjectRecord{}
	}

	res.Total = total
	res.Offset = *offset
	res.Limit = *limit

	return res, nil
}

type activityRecordListItem struct {
	RecordId            string     `db:"record_id"`
	UserId              string     `db:"user_id"`
	Username            string     `db:"username"`
	Place               string     `db:"place"`
	CheckedInAt         time.Time  `db:"checked_in_at"`
	CheckedOutAt        *time.Time `db:"checked_out_at"`
	InitialCheckedInAt  time.Time  `db:"initial_checked_in_at"`
	InitialCheckedOutAt *time.Time `db:"initial_checked_out_at"`
}

func selectActivityRecords(dbClient db.Client, offset int, limit int) ([]activityRecordListItem, *response.Error) {
	params := struct {
		Offset int `twowaysql:"offset"`
		Limit  int `twowaysql:"limit"`
	}{
		Offset: offset,
		Limit:  limit,
	}

	records := []activityRecordListItem{}
	err := dbClient.Select(&records, "sql/activity/select_activity_records.sql", &params)
	if err != nil {
		return []activityRecordListItem{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "アクティビティレコード一覧の取得に失敗しました",
			Log:     err.Error(),
		}
	}
	return records, nil
}

type activityRecordsCount struct {
	Total int `db:"total"`
}

func selectActivityRecordsCount(dbClient db.Client) (int, *response.Error) {
	counts := []activityRecordsCount{}
	err := dbClient.Select(&counts, "sql/activity/select_activity_records_count.sql", nil)
	if err != nil {
		return 0, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "アクティビティレコード総数の取得に失敗しました",
			Log:     err.Error(),
		}
	}
	if len(counts) == 0 {
		return 0, nil
	}
	return counts[0].Total, nil
}
