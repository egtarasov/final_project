package student_client

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/student_repo"
)

func (s *StudentClient) Delete(ctx context.Context, id int64) (bool, error) {
	tp := otel.Tracer("ClientDelete")
	ctx, span := tp.Start(ctx, "start deleting")
	span.SetAttributes(
		attribute.Key("param_student").Int64(id))
	defer span.End()

	response, err := s.grpc.DeleteStudentById(ctx, &student_repo.DeleteStudentRequest{Id: id})
	if err != nil {
		return false, err
	}
	return response.Ok, nil
}
