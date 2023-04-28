package student

import (
	"database/sql"
	"time"
)

type Student struct {
	Id             int64          `db:"id"`
	FirstName      string         `db:"fisrt_name"`
	SecondName     string         `db:"second_name"`
	MiddleName     sql.NullString `db:"middle_name"`
	Gpa            float64        `db:"gpa"`
	AttendanceRate float64        `db:"attendance_rate"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	GroupId        sql.NullInt64  `db:"group_id"`
}
