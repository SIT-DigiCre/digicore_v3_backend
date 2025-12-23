package activity

import "time"

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
