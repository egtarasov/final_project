package student

import (
	"context"
	"database/sql"
	"errors"
	"homework-5/internal/app/database"
)

var ErrObjectNotFound = errors.New("Object not found")

type StudentsRepository struct {
	client database.Dbops
}

func NewStudentsRepository(client database.Dbops) *StudentsRepository {
	return &StudentsRepository{client: client}
}

func (s *StudentsRepository) GetById(ctx context.Context, id int64) (*Student, error) {
	var student Student
	query := `SELECT id, fisrt_name, second_name, middle_name, gpa, attendance_rate, created_at, updated_at, group_id from students WHERE id = $1`
	err := s.client.Get(ctx, &student, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrObjectNotFound
	}

	return &student, err
}

func (s *StudentsRepository) Add(ctx context.Context, student *Student) (int64, error) {
	query := `INSERT INTO students (fisrt_name, second_name, middle_name, gpa, attendance_rate, group_id) 
				 VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int64

	err := s.client.ExecQueryRow(ctx, query,
		student.FirstName,
		student.SecondName,
		student.MiddleName,
		student.Gpa,
		student.AttendanceRate,
		student.GroupId).Scan(&id)

	return id, err
}

func (s *StudentsRepository) UpdateById(ctx context.Context, id int64, student *Student) (bool, error) {
	query := `UPDATE students
			  SET fisrt_name = $1, second_name = $2, middle_name = $3, gpa = $4, attendance_rate = $5, created_at = $6, updated_at = $7, group_id = $8
			  WHERE id = $9;`

	result, err := s.client.Exec(ctx, query,
		student.FirstName,
		student.SecondName,
		student.MiddleName,
		student.Gpa,
		student.AttendanceRate,
		student.CreatedAt,
		student.UpdatedAt,
		student.GroupId,
		id)

	return result.RowsAffected() > 0, err
}

func (s *StudentsRepository) Remove(ctx context.Context, id int64) (bool, error) {
	result, err := s.client.Exec(ctx, "DELETE FROM students WHERE id = $1", id)

	return result.RowsAffected() > 0, err
}
