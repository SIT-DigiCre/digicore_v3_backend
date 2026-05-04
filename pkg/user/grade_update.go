package user

// 学年補正申請の最大承認回数
const maxApprovedCount = 2

// DB構造体: ユーザーの申請一覧用
type gradeUpdate struct {
	GradeUpdateId string `db:"grade_update_id"`
	GradeDiff     int    `db:"grade_diff"`
	Reason        string `db:"reason"`
	Status        string `db:"status"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
}

// DB構造体: 承認済み回数カウント用
type approvedCount struct {
	ApprovedCount int `db:"approved_count"`
}

// DB構造体: 未処理件数カウント用
type pendingCount struct {
	PendingCount int `db:"pending_count"`
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
