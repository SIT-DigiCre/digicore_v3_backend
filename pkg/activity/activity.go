package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

type ActivityRecord struct {
	ID                  string     `db:"id"`
	UserID              string     `db:"user_id"`
	Place               string     `db:"place"`
	Note                *string    `db:"note"`
	InitialCheckedInAt  time.Time  `db:"initial_checked_in_at"`
	InitialCheckedOutAt *time.Time `db:"initial_checked_out_at"`
	CheckedInAt         time.Time  `db:"checked_in_at"`
	CheckedOutAt        *time.Time `db:"checked_out_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

// アクティビティ入退室時刻更新リクエスト
type ActivityRecordUpdateRequest struct {
	ActivityType string    `json:"activity_type" validate:"required,oneof=checkin checkout" ja:"アクティビティタイプ"`
	Time         time.Time `json:"time" validate:"required" ja:"更新時刻"`
}

func selectLatestActivity(dbClient db.TransactionClient, userId string, place string) (*ActivityRecord, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
		Place  string `twowaysql:"place"`
	}{
		UserId: userId,
		Place:  place,
	}

	records := []ActivityRecord{}
	err := dbClient.Select(&records, "sql/activity/select_latest_activity.sql", &params)
	if err != nil {
		return nil, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}

	if len(records) == 0 {
		return nil, nil
	}

	return &records[0], nil
}
