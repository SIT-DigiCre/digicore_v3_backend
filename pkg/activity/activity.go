package activity

import "time"

// アクティビティレコードの種別
const (
	TypeCheckIn  = "checkin"
	TypeCheckOut = "checkout"
)

// ActivityRecord は activity_records テーブルの1レコードを表します。
// 将来的にチェックイン/チェックアウトのユースケースや在室判定ロジックで利用されます。
type ActivityRecord struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Place     string    `db:"place"`
	Type      string    `db:"type"`
	Datetime  time.Time `db:"datetime"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}


