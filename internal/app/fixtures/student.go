package fixtures

import (
	"database/sql"
	"homework-5/internal/app/student"
	"time"
)

const (
	Id             = 0
	Firstname      = "Obu"
	SecondName     = "Obu"
	MiddleName     = "Obu"
	Gpa            = 9.3
	AttendanceRate = 0.3
	GroupId        = 12
)

var (
	Time, _ = time.Parse("2006-01-02", "2023-04-21")
)

type StudentBuilder struct {
	student student.Student
}

func (s *StudentBuilder) Id(id int64) *StudentBuilder {
	s.student.Id = id
	return s
}

func (s *StudentBuilder) FirstName(firstName string) *StudentBuilder {
	s.student.FirstName = firstName
	return s
}

func (s *StudentBuilder) SecondName(secondName string) *StudentBuilder {
	s.student.SecondName = secondName
	return s
}

func (s *StudentBuilder) MiddleName(middleName string) *StudentBuilder {
	s.student.MiddleName = sql.NullString{String: middleName, Valid: true}
	return s
}

func (s *StudentBuilder) AttendanceRate(attendanceRate float64) *StudentBuilder {
	s.student.AttendanceRate = attendanceRate
	return s
}

func (s *StudentBuilder) Gpa(gpa float64) *StudentBuilder {
	s.student.Gpa = gpa
	return s
}

func (s *StudentBuilder) CreatedAt(time time.Time) *StudentBuilder {
	s.student.CreatedAt = time
	return s
}

func (s *StudentBuilder) GroupId(GroupId int64) *StudentBuilder {
	s.student.GroupId = sql.NullInt64{Int64: GroupId, Valid: true}
	return s
}

func (s *StudentBuilder) UpdatedAt(time time.Time) *StudentBuilder {
	s.student.UpdatedAt = sql.NullTime{Time: time, Valid: true}
	return s
}

func (s *StudentBuilder) P() *student.Student {
	return &s.student
}

func (s *StudentBuilder) V() student.Student {
	return s.student
}

func DefaultStudent() *StudentBuilder {
	student := StudentBuilder{}
	return student.
		Id(Id).
		FirstName(Firstname).
		SecondName(SecondName).
		MiddleName(MiddleName).
		GroupId(GroupId).
		Gpa(Gpa).
		CreatedAt(Time).
		AttendanceRate(AttendanceRate)
}
