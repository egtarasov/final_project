package core

import (
	"context"
	"database/sql"
	"fmt"
	"homework-5/internal/app/group"
	"strconv"
)

const (
	getParamsLenGroups    = 1
	deleteParamsLenGroups = 1
	addParamsLenGroups    = 3
	updateParamsLenGroups = 3
)

type groupCommand struct {
	groupRepo GroupRepository
	response  string
}

func NewGroupCommand(groupRepo GroupRepository) *groupCommand {
	return &groupCommand{
		groupRepo: groupRepo,
	}
}

func (g *groupCommand) Process(ctx context.Context, params []string) (string, error) {
	switch params[0] {
	case "Get":
		return g.getCommandGroup(ctx, params)
	case "Add":
		return g.addCommandGroup(ctx, params)
	case "Update":
		return g.updateCommandGroup(ctx, params)
	case "Delete":
		return g.deleteCommandGroup(ctx, params)
	default:
		return "", InvalidInput
	}
}

func (g *groupCommand) deleteCommandGroup(ctx context.Context, params []string) (string, error) {
	if len(params[1:]) != deleteParamsLenGroups {
		return "", InvalidInput
	}
	id, err := getId(params)
	if err != nil {
		return "", InvalidInput
	}
	ok, err := g.groupRepo.Remove(ctx, id)
	if err != nil {
		return "", ProcessingError
	}
	if ok {
		g.response = fmt.Sprintf("Group  with id[%v] has been deleted\n", id)
	} else {
		g.response = fmt.Sprintf("Cant delete group with id[%v]", id)
	}
	return g.response, nil
}

func (g *groupCommand) updateCommandGroup(ctx context.Context, params []string) (string, error) {
	if len(params[1:]) != updateParamsLenGroups {
		return "", InvalidInput
	}
	group, err := g.getGroup(params[1:])
	if err != nil {
		return "", InvalidInput
	}
	ok, err := g.groupRepo.UpdateById(ctx, uint64(group.Id), group)
	if err != nil {
		return "", ProcessingError
	}
	if ok {
		g.response = fmt.Sprintf("Group with id[%v] has been updated\n", group.Id)
	} else {
		g.response = fmt.Sprintf("Cant update group with id[%v]", group.Id)
	}
	return g.response, nil
}

func (g *groupCommand) addCommandGroup(ctx context.Context, params []string) (string, error) {
	if len(params[1:]) != addParamsLenGroups {
		return "", InvalidInput
	}
	group, err := g.getGroup(params[1:])
	if err != nil {
		return "", InvalidInput
	}
	id, err := g.groupRepo.Add(ctx, group)
	if err != nil {
		return "", ProcessingError
	}
	g.response = fmt.Sprintf("Group has been added with id[%v]\n", id)
	return g.response, nil
}

func (g *groupCommand) getGroup(params []string) (*group.Group, error) {
	id, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return nil, InvalidInput
	}

	valid := true
	if params[1] == "null" {
		valid = false
	}

	year, err := strconv.ParseUint(params[2], 10, 32)
	return &group.Group{
		Id:   id,
		Name: sql.NullString{String: params[1], Valid: valid},
		Year: int32(year),
	}, nil
}

func (g *groupCommand) getCommandGroup(ctx context.Context, params []string) (string, error) {
	if len(params[1:]) != getParamsLenGroups {
		return "", InvalidInput
	}
	id, err := getId(params)
	if err != nil {
		return "", InvalidInput
	}
	group, err := g.groupRepo.GetById(ctx, id)
	if err != nil {
		return "", ProcessingError
	}
	g.response = fmt.Sprintf("group: %v\n", group)
	return g.response, nil
}
