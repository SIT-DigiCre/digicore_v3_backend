package admin

type schoolGradeUpdateTarget struct {
	UserId             string `db:"user_id"`
	StudentNumber      string `db:"student_number"`
	ApprovedGradeDiffs int    `db:"approved_grade_diffs"`
}
