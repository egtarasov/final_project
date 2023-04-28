package group_client

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/pb/group_repo"
)

func (s *StudentClient) Update(ctx context.Context, id int64, group *group_repo.Group) (bool, error) {
	tp := otel.Tracer("ClientUpdate_Group")
	ctx, span := tp.Start(ctx, "start updating group")
	marshalled, err := json.Marshal(group)
	if err != nil {
		marshalled = []byte("cant marshall group")
	}
	span.SetAttributes(
		attribute.Key("param_group").String(string(marshalled)),
		attribute.Key("param_id").Int64(id))
	defer span.End()

	response, err := s.grpc.UpdateGroup(ctx, &group_repo.UpdateGroupRequest{
		Group: group,
		Id:    id,
	})
	if err != nil {
		return false, err
	}
	return response.Ok, nil
}
