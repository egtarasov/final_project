package student_client

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/student_repo"
)

func (s *StudentClient) Update(ctx context.Context, id int64, st *student_repo.Student) (bool, error) {
	tp := otel.Tracer("ClientUpdate")
	ctx, span := tp.Start(ctx, "start updating")
	marshalled, err := json.Marshal(st)
	if err != nil {
		marshalled = []byte("cant marshall struct")
	}
	span.SetAttributes(
		attribute.Key("param_student").String(string(marshalled)),
		attribute.Key("param_id").Int64(id))
	defer span.End()

	response, err := s.grpc.UpdateStudent(ctx, &student_repo.UpdateStudentRequest{
		Student: st,
		Id:      id,
	})
	if err != nil {
		return false, err
	}
	return response.Ok, nil
}
