package group_client

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/group_repo"
)

func (s *GroupClient) GetById(ctx context.Context, id int64) (*group_repo.Group, error) {
	tp := otel.Tracer("ClientGet_Group")
	ctx, span := tp.Start(ctx, "start retrieving group")

	span.SetAttributes(
		attribute.Key("param_id").Int64(id))
	defer span.End()
	response, err := s.grpc.GetGroupById(ctx, &group_repo.GetGroupRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return response.Group, nil
}
