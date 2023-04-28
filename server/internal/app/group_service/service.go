package group_service

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/internal/app"
	"homework-5/internal/app/group"
	"homework-5/internal/app/pb/group_repo"
)

type Implementation struct {
	group_repo.UnsafeGroupServiceServer

	repo group.GroupRepository
}

func NewGroupService(repo group.GroupRepository) *Implementation {
	return &Implementation{repo: repo}
}

func (i *Implementation) CreateGroup(ctx context.Context, request *group_repo.CreateGroupRequest) (*group_repo.CreateGroupResponse, error) {
	app.GroupOpProcessed.Inc()

	tr := otel.Tracer("CreateGroup")
	ctx, span := tr.Start(ctx, "received request")
	marshalled, err := json.Marshal(request.Group)
	if err != nil {
		return nil, err
	}
	span.SetAttributes(attribute.Key("params_group").String(string(marshalled)))
	defer span.End()

	id, err := i.repo.Add(ctx, ParseGroup(request.Group))
	if err != nil {
		return nil, fmt.Errorf("cant create group")
	}

	return &group_repo.CreateGroupResponse{Id: id}, nil
}

func (i *Implementation) GetGroupById(ctx context.Context, request *group_repo.GetGroupRequest) (*group_repo.GetGroupResponse, error) {
	app.GroupOpProcessed.Inc()

	tr := otel.Tracer("GetGroup")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params_id").Int64(request.Id))
	defer span.End()

	group, err := i.repo.GetById(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("cant create group")
	}

	return &group_repo.GetGroupResponse{Group: ParseGroupRequest(group)}, nil
}

func (i *Implementation) DeleteGroupById(ctx context.Context, request *group_repo.DeleteGroupRequest) (*group_repo.DeleteGroupResponse, error) {
	app.GroupOpProcessed.Inc()

	tr := otel.Tracer("DeleteGroup")
	ctx, span := tr.Start(ctx, "received request")
	span.SetAttributes(attribute.Key("params_id").Int64(request.Id))
	defer span.End()

	ok, err := i.repo.Remove(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("cant create group")
	}

	return &group_repo.DeleteGroupResponse{Ok: ok}, nil
}

func (i *Implementation) UpdateGroup(ctx context.Context, request *group_repo.UpdateGroupRequest) (*group_repo.UpdateGroupResponse, error) {
	app.GroupOpProcessed.Inc()

	tr := otel.Tracer("UpdateGroup")
	ctx, span := tr.Start(ctx, "received request")
	marshalled, err := json.Marshal(request.Group)
	if err != nil {
		return nil, err
	}
	span.SetAttributes(attribute.Key("params_group").String(string(marshalled)))
	span.SetAttributes(attribute.Key("params_id").Int64(request.Id))
	defer span.End()

	ok, err := i.repo.UpdateById(ctx, request.Id, ParseGroup(request.Group))
	if err != nil {
		return nil, fmt.Errorf("cant create group")
	}

	return &group_repo.UpdateGroupResponse{Ok: ok}, nil
}
