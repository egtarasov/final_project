package group

import "database/sql"

type Group struct {
	Id   int64          `db:"id"`
	Name sql.NullString `db:"group_name"`
	Year int32          `db:"st_year"`
}
