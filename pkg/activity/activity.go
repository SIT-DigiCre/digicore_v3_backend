package activity

import "time"

type ActivityRecord struct {
	ID                string     `db:"id"`
	UserID            string     `db:"user_id"`
	Place             string     `db:"place"`
	Note              *string    `db:"note"`
	InitialCheckInAt  time.Time  `db:"initial_check_in_at"`
	InitialCheckOutAt *time.Time `db:"initial_check_out_at"`
	CheckInAt         time.Time  `db:"check_in_at"`
	CheckOutAt        *time.Time `db:"check_out_at"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}
