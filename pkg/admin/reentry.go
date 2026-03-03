package admin

// DB構造体: 管理者向け再入部申請一覧用
type adminReentry struct {
	ReentryId     string `db:"reentry_id"`
	UserId        string `db:"user_id"`
	Username      string `db:"username"`
	StudentNumber string `db:"student_number"`
	ReentryCount  int    `db:"reentry_count"`
	Status        string `db:"status"`
	Note          string `db:"note"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
}

type reentryDetail struct {
	ReentryId    string `db:"reentry_id"`
	UserId       string `db:"user_id"`
	ReentryCount int    `db:"reentry_count"`
	Status       string `db:"status"`
	Note         string `db:"note"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}
