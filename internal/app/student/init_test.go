package student

import (
	"context"
	"database/sql"
	"fmt"
	"homework-5/internal/app/database"
	"strings"
	"sync"
	"testing"
)

type TestDb struct {
	sync.Mutex
	ctx  context.Context
	db   *database.Database
	repo *StudentsRepository
}

func NewTestDb() *TestDb {
	ctx := context.Background()
	db, err := database.NewDb(ctx)
	if err != nil {
		return nil
	}
	return &TestDb{
		ctx:  ctx,
		repo: NewStudentsRepository(db),
		db:   db,
	}
}

func (db *TestDb) setUp(t *testing.T) {
	db.Lock()
	db.truncate()
}

func (db *TestDb) tearDown(t *testing.T) {
	defer db.Unlock()
	db.truncate()
}

func (db *TestDb) truncate() {
	var tables []string
	err := db.db.Select(db.ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'goose_db_version'")
	if err != nil {
		return
	}
	if len(tables) == 0 {
		panic("not tables in db")
	}
	query := fmt.Sprintf("TRUNCATE TABLE %v", strings.Join(tables, ", "))
	if _, err := db.db.Exec(db.ctx, query); err != nil {
		panic(err)
	}
}

func (db *TestDb) add(student *Student, id *uint64) {
	var groupId int64
	db.db.ExecQueryRow(db.ctx, `INSERT INTO groups (group_name, st_year) VALUES('test_group', 3) RETURNING id`).Scan(&groupId)

	student.GroupId = sql.NullInt64{Int64: groupId, Valid: true}
	db.db.ExecQueryRow(
		db.ctx,
		`INSERT INTO students (fisrt_name, second_name, middle_name, gpa, attendance_rate, created_at, updated_at, group_id) 
				 VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		student.FirstName,
		student.SecondName,
		student.MiddleName,
		student.Gpa,
		student.AttendanceRate,
		student.CreatedAt,
		student.UpdatedAt,
		student.GroupId).Scan(id)
}
