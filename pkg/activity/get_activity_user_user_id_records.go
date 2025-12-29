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

func GetActivityUserUserIdRecords(ctx echo.Context, dbClient db.Client, userId string, place *string, offset *int, limit *int) (api.ResGetActivityUserUserIdRecords, *response.Error) {
	res := api.ResGetActivityUserUserIdRecords{}

	// デフォルト値の設定
	defaultOffset := 0
	defaultLimit := 50
	if offset == nil {
		offset = &defaultOffset
	}
	if limit == nil {
		limit = &defaultLimit
	}

	records, err := selectUserRecords(dbClient, userId, place, *offset, *limit)
	if err != nil {
		return api.ResGetActivityUserUserIdRecords{}, err
	}

	total, err := selectUserRecordsCount(dbClient, userId, place)
	if err != nil {
		return api.ResGetActivityUserUserIdRecords{}, err
	}

	rerr := copier.Copy(&res.Records, &records)
	if rerr != nil {
		return api.ResGetActivityUserUserIdRecords{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "入室記録の取得に失敗しました",
			Log:     rerr.Error(),
		}
	}
	if res.Records == nil {
		res.Records = []api.ResGetActivityUserUserIdRecordsObjectRecord{}
	}

	res.Total = total
	res.Offset = *offset
	res.Limit = *limit

	return res, nil
}

type userRecord struct {
	RecordId            string     `db:"record_id"`
	Place               string     `db:"place"`
	CheckedInAt         time.Time  `db:"checked_in_at"`
	CheckedOutAt        *time.Time `db:"checked_out_at"`
	InitialCheckedInAt  time.Time  `db:"initial_checked_in_at"`
	InitialCheckedOutAt *time.Time `db:"initial_checked_out_at"`
}

func selectUserRecords(dbClient db.Client, userId string, place *string, offset int, limit int) ([]userRecord, *response.Error) {
	params := struct {
		UserId string  `twowaysql:"userId"`
		Place  *string `twowaysql:"place"`
		Offset int     `twowaysql:"offset"`
		Limit  int     `twowaysql:"limit"`
	}{
		UserId: userId,
		Place:  place,
		Offset: offset,
		Limit:  limit,
	}
	records := []userRecord{}
	err := dbClient.Select(&records, "sql/activity/select_user_records.sql", &params)
	if err != nil {
		return []userRecord{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "入室記録の取得に失敗しました",
			Log:     err.Error(),
		}
	}
	return records, nil
}

type userRecordsCount struct {
	Total int `db:"total"`
}

func selectUserRecordsCount(dbClient db.Client, userId string, place *string) (int, *response.Error) {
	params := struct {
		UserId string  `twowaysql:"userId"`
		Place  *string `twowaysql:"place"`
	}{
		UserId: userId,
		Place:  place,
	}
	counts := []userRecordsCount{}
	err := dbClient.Select(&counts, "sql/activity/select_user_records_count.sql", &params)
	if err != nil {
		return 0, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "入室記録数の取得に失敗しました",
			Log:     err.Error(),
		}
	}
	if len(counts) == 0 {
		return 0, nil
	}
	return counts[0].Total, nil
}
