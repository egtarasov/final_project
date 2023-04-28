package group_client

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/group_repo"
)

func (s *StudentClient) Create(ctx context.Context, group *group_repo.Group) (int64, error) {
	tp := otel.Tracer("ClientCreate_Group")
	ctx, span := tp.Start(ctx, "start creating group")
	marshalled, err := json.Marshal(group)
	if err != nil {
		marshalled = []byte("cant marshall group")
	}
	span.SetAttributes(
		attribute.Key("param_group").String(string(marshalled)))
	defer span.End()

	response, err := s.grpc.CreateGroup(ctx, &group_repo.CreateGroupRequest{
		Group: group,
	})
	if err != nil {
		return 0, err
	}
	return response.Id, nil
}
