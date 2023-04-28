package student_service

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/internal/app"
	"homework-5/internal/app/pb/student_repo"
	"homework-5/internal/app/student"
)

type Implementation struct {
	student_repo.UnsafeStudentServiceServer
	repo student.Repository
}

func NewImplementation(repo student.Repository) *Implementation {
	return &Implementation{repo: repo}
}

func (i *Implementation) CreateStudent(ctx context.Context, request *student_repo.CreateStudentRequest) (*student_repo.CreateStudentResponse, error) {
	app.StudentOpProcessed.Inc()

	tr := otel.Tracer("CreateStudent")
	ctx, span := tr.Start(ctx, "received request")
	marshalled, err := json.Marshal(request.Student)
	if err != nil {
		return nil, err
	}
	span.SetAttributes(attribute.Key("params_student").String(string(marshalled)))
	defer span.End()

	id, err := i.repo.Add(ctx, ParseStudent(request.Student))
	if err != nil {
		return nil, fmt.Errorf("cant create student")
	}
	return &student_repo.CreateStudentResponse{Id: id}, nil
}

func (i *Implementation) GetStudentById(ctx context.Context, request *student_repo.GetStudentRequest) (*student_repo.GetStudentResponse, error) {
	app.StudentOpProcessed.Inc()

	tr := otel.Tracer("GetStudent")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params_id").Int64(request.Id))
	defer span.End()

	st, err := i.repo.GetById(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("cant get student")
	}
	return &student_repo.GetStudentResponse{Student: ParseStudentRequest(st)}, nil
}

func (i *Implementation) DeleteStudentById(ctx context.Context, request *student_repo.DeleteStudentRequest) (*student_repo.DeleteStudentResponse, error) {
	app.StudentOpProcessed.Inc()

	tr := otel.Tracer("DeleteStudent")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params_id").Int64(request.Id))
	defer span.End()

	ok, err := i.repo.Remove(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("cant remove student")
	}
	return &student_repo.DeleteStudentResponse{Ok: ok}, nil
}

func (i *Implementation) UpdateStudent(ctx context.Context, request *student_repo.UpdateStudentRequest) (*student_repo.UpdateStudentResponse, error) {
	app.StudentOpProcessed.Inc()

	tr := otel.Tracer("UpdateStudent")
	ctx, span := tr.Start(ctx, "received request")
	marshalled, err := json.Marshal(request.Student)
	if err != nil {
		return nil, err
	}
	span.SetAttributes(attribute.Key("params_student").String(string(marshalled)))
	span.SetAttributes(attribute.Key("params_id").Int64(request.Id))
	defer span.End()

	ok, err := i.repo.UpdateById(ctx, request.Id, ParseStudent(request.Student))
	if err != nil {
		return nil, fmt.Errorf("cant remove student")
	}
	return &student_repo.UpdateStudentResponse{Ok: ok}, nil
}
