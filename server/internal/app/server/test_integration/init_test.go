package test_integration

import (
	"context"
	"database/sql"
	"fmt"
	database2 "homework-5/server/internal/app/database"
	group2 "homework-5/server/internal/app/group"
	"homework-5/server/internal/app/server"
	student2 "homework-5/server/internal/app/student"
	"strings"
	"sync"
)

type TestHandler struct {
	sync.Mutex
	ctx    context.Context
	server *server.Server
	db     *database2.Database
}

func NewTestHandler() *TestHandler {
	ctx := context.Background()
	db, _ := database2.NewDb(ctx)
	groupRepo := group2.NewGroupsRepository(db)
	studentRepo := student2.NewStudentsRepository(db)
	server := server.NewServer(ctx, studentRepo, groupRepo)
	return &TestHandler{
		ctx:    ctx,
		server: server,
		db:     db,
	}
}

func (t *TestHandler) setUp(student *student2.Student, group *group2.Group) {
	t.Lock()
	t.truncate()
	t.add(student, group)
}

func (t *TestHandler) tearDown() {
	defer t.Unlock()
	t.truncate()
}

func (t *TestHandler) truncate() {
	var tables []string
	err := t.db.Select(t.ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'goose_db_version'")
	if err != nil {
		return
	}
	if len(tables) == 0 {
		panic("not tables in db")
	}
	query := fmt.Sprintf("TRUNCATE TABLE %v", strings.Join(tables, ", "))
	if _, err := t.db.Exec(t.ctx, query); err != nil {
		panic(err)
	}
}

func (t *TestHandler) add(student *student2.Student, group *group2.Group) {
	t.db.ExecQueryRow(t.ctx, `INSERT INTO groups (group_name, st_year) VALUES($1, $2) RETURNING id`,
		group.Name, group.Year).Scan(&group.Id)

	student.GroupId = sql.NullInt64{Int64: group.Id, Valid: true}
	t.db.ExecQueryRow(
		t.ctx,
		`INSERT INTO students (fisrt_name, second_name, middle_name, gpa, attendance_rate, created_at, updated_at, group_id) 
				 VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		student.FirstName,
		student.SecondName,
		student.MiddleName,
		student.Gpa,
		student.AttendanceRate,
		student.CreatedAt,
		student.UpdatedAt,
		student.GroupId).Scan(&student.Id)
}
