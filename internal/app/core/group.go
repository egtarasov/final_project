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
	addParamsLenGroups    = 4
	updateParamsLenGroups = 4
)

type groupCommand struct {
	groupRepo GroupRepository
}

func NewGroupCommand(groupRepo GroupRepository) *groupCommand {
	return &groupCommand{
		groupRepo: groupRepo,
	}
}

func (g *groupCommand) Process(ctx context.Context, params []string) error {
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
		return InvalidInput
	}
}

func (g *groupCommand) deleteCommandGroup(ctx context.Context, params []string) error {
	if len(params[1:]) != deleteParamsLenGroups {
		return InvalidInput
	}
	id, err := getId(params)
	if err != nil {
		return InvalidInput
	}
	ok, err := g.groupRepo.Remove(ctx, id)
	if err != nil {
		return ProcessingError
	}
	if ok {
		fmt.Printf("Group  with id[%v] has been deleted\n", id)
	} else {
		fmt.Printf("Cant delete group with id[%v]", id)
	}
	return nil
}

func (g *groupCommand) updateCommandGroup(ctx context.Context, params []string) error {
	if len(params[1:]) != updateParamsLenGroups {
		return InvalidInput
	}
	group, err := g.getGroup(params[1:])
	if err != nil {
		return InvalidInput
	}
	ok, err := g.groupRepo.UpdateById(ctx, uint64(group.Id), group)
	if err != nil {
		return ProcessingError
	}
	if ok {
		fmt.Printf("Group with id[%v] has been updated\n", group.Id)
	} else {
		fmt.Printf("Cant update group with id[%v]", group.Id)
	}
	return nil
}

func (g *groupCommand) addCommandGroup(ctx context.Context, params []string) error {
	if len(params[1:]) != addParamsLenGroups {
		return InvalidInput
	}
	group, err := g.getGroup(params[1:])
	if err != nil {
		return InvalidInput
	}
	id, err := g.groupRepo.Add(ctx, group)
	if err != nil {
		return ProcessingError
	}
	fmt.Printf("Group has been added with id[%v]\n", id)
	return nil
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

	amountOfStudents, err := strconv.ParseInt(params[2], 10, 64)
	if err != nil {
		return nil, InvalidInput
	}

	year, err := strconv.ParseUint(params[3], 10, 32)
	return &group.Group{
		Id:               id,
		Name:             sql.NullString{String: params[1], Valid: valid},
		AmountOfStudents: amountOfStudents,
		Year:             int32(year),
	}, nil
}

func (g *groupCommand) getCommandGroup(ctx context.Context, params []string) error {
	if len(params[1:]) != getParamsLenGroups {
		return InvalidInput
	}
	id, err := getId(params)
	if err != nil {
		return InvalidInput
	}
	group, err := g.groupRepo.GetById(ctx, id)
	if err != nil {
		return ProcessingError
	}
	fmt.Printf("group: %v\n", group)
	return nil
}
