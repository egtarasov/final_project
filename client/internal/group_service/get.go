package student_client

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/student_repo"
)

func (s *StudentClient) GetById(ctx context.Context, id int64) (*student_repo.Student, error) {
	tp := otel.Tracer("ClientGet")
	ctx, span := tp.Start(ctx, "start retrieving")

	span.SetAttributes(
		attribute.Key("param_id").Int64(id))
	defer span.End()
	response, err := s.grpc.GetStudentById(ctx, &student_repo.GetStudentRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return response.Student, nil
}
