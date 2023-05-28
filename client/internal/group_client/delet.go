package group_client

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/group_repo"
)

func (s *GroupClient) Delete(ctx context.Context, id int64) (bool, error) {
	tp := otel.Tracer("ClientDelete_Group")
	ctx, span := tp.Start(ctx, "start deleting group")
	span.SetAttributes(
		attribute.Key("param_group").Int64(id))
	defer span.End()

	response, err := s.grpc.DeleteGroupById(ctx, &group_repo.DeleteGroupRequest{Id: id})
	if err != nil {
		return false, err
	}
	return response.Ok, nil
}
