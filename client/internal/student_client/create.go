package student_client

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/student_repo"
)

func (s *StudentClient) Create(ctx context.Context, st *student_repo.Student) (int64, error) {
	tp := otel.Tracer("ClientCreate")
	ctx, span := tp.Start(ctx, "start creating")
	marshalled, err := json.Marshal(st)
	if err != nil {
		marshalled = []byte("cant marshall struct")
	}
	span.SetAttributes(
		attribute.Key("param_student").String(string(marshalled)))
	defer span.End()

	response, err := s.grpc.CreateStudent(ctx, &student_repo.CreateStudentRequest{
		Student: st,
	})
	if err != nil {
		return 0, err
	}
	return response.Id, nil
}
