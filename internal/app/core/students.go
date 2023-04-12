package core

import (
	"context"
	"database/sql"
	"fmt"
	"homework-5/internal/app/student"
	"strconv"
	"time"
)

const (
	getParamsLenStudents    = 1
	deleteParamsLenStudents = 1
	addParamsLenStudents    = 7
	updateParamsLenStudents = 7
)

type studentCommand struct {
	repo StudentRepository
}

func NewStudentCommand(repo StudentRepository) *studentCommand {
	return &studentCommand{
		repo: repo,
	}
}

func (s *studentCommand) Process(ctx context.Context, params []string) error {
	switch params[0] {
	case "Get":
		return s.getCommandStudent(ctx, params)
	case "Add":
		return s.addCommandStudent(ctx, params)
	case "Update":
		return s.updateCommandStudent(ctx, params)
	case "Delete":
		return s.deleteCommandStudent(ctx, params)
	default:
		return InvalidInput
	}
}

func (s *studentCommand) deleteCommandStudent(ctx context.Context, params []string) error {
	if len(params[1:]) != deleteParamsLenStudents {
		return InvalidInput
	}
	id, err := getId(params)
	if err != nil {
		return InvalidInput
	}
	ok, err := s.repo.Remove(ctx, id)
	if err != nil {
		return ProcessingError
	}
	if ok {
		fmt.Printf("Student has been added with id[%v]\n", id)
	} else {
		fmt.Printf("Cant update student with id[%v]", id)
	}
	return nil
}

func (s *studentCommand) updateCommandStudent(ctx context.Context, params []string) error {
	if len(params[1:]) != updateParamsLenStudents {
		return InvalidInput
	}
	student, err := s.getStudent(params[1:])
	if err != nil {
		return InvalidInput
	}
	student.UpdatedAt.Valid = true
	ok, err := s.repo.UpdateById(ctx, uint64(student.Id), student)
	if err != nil {
		return ProcessingError
	}
	if ok {
		fmt.Printf("Student has been added with id[%v]\n", student.Id)
	} else {
		fmt.Printf("Cant update student with id[%v]", student.Id)
	}
	return nil
}

func (s *studentCommand) addCommandStudent(ctx context.Context, params []string) error {
	if len(params[1:]) != addParamsLenStudents {
		return InvalidInput
	}
	student, err := s.getStudent(params[1:])
	if err != nil {
		return InvalidInput
	}
	id, err := s.repo.Add(ctx, student)
	if err != nil {
		return ProcessingError
	}
	fmt.Printf("Student has been added with id[%v]\n", id)
	return nil
}

func (s *studentCommand) getStudent(params []string) (*student.Student, error) {
	id, err := strconv.ParseInt(params[0], 01, 64)
	if err != nil {
		return nil, InvalidInput
	}

	gpa, err := strconv.ParseFloat(params[4], 64)
	if err != nil {
		return nil, InvalidInput
	}

	attendanceRate, err := strconv.ParseFloat(params[5], 64)
	if err != nil {
		return nil, InvalidInput
	}

	var groupId int64
	valid := true
	if params[8] == "null" {
		valid = false
	}
	if valid {
		groupId, err = strconv.ParseInt(params[8], 10, 64)
		if err != nil {
			return nil, InvalidInput
		}
	}
	return &student.Student{
		Id:             id,
		FirstName:      params[1],
		SecondName:     params[2],
		MiddleName:     sql.NullString{String: params[3], Valid: false},
		Gpa:            gpa,
		AttendanceRate: attendanceRate,
		CreatedAt:      time.Now(), //default
		UpdatedAt:      sql.NullTime{Time: time.Now(), Valid: false},
		GroupId:        sql.NullInt64{Int64: groupId, Valid: valid},
	}, nil
}

func (s *studentCommand) getCommandStudent(ctx context.Context, params []string) error {
	fmt.Println(params)
	if len(params[1:]) != getParamsLenStudents {
		return InvalidInput
	}
	id, err := getId(params)
	if err != nil {
		return InvalidInput
	}
	student, err := s.repo.GetById(ctx, id)
	if err != nil {
		return ProcessingError
	}
	fmt.Printf("student: %v\n", student)
	return nil
}
