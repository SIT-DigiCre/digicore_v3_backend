package activity

import "time"

// アクティビティレコードの種別
const (
	TypeCheckIn  = "checkin"
	TypeCheckOut = "checkout"
)

type ActivityRecord struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	Place        string    `db:"place"`
	ActivityType string    `db:"activity_type"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
