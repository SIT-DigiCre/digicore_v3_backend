package admin

// DB構造体: 管理者向け申請一覧用（ユーザー情報付き）
type adminGradeUpdate struct {
	GradeUpdateId string `db:"grade_update_id"`
	UserId        string `db:"user_id"`
	Username      string `db:"username"`
	GradeDiff     int    `db:"grade_diff"`
	Reason        string `db:"reason"`
	Status        string `db:"status"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
}

// DB構造体: ID指定取得用
type gradeUpdateDetail struct {
	GradeUpdateId string `db:"grade_update_id"`
	UserId        string `db:"user_id"`
	GradeDiff     int    `db:"grade_diff"`
	Reason        string `db:"reason"`
	Status        string `db:"status"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
}
