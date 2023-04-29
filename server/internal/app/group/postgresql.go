package group

import (
	"context"
	"database/sql"
	"errors"
	"homework-5/server/internal/app/database"
)

var ErrObjectNotFound = errors.New("Object not found")

type GroupsRepository struct {
	client database.Dbops
}

func NewGroupsRepository(client database.Dbops) *GroupsRepository {
	return &GroupsRepository{client: client}
}

func (s *GroupsRepository) GetById(ctx context.Context, id int64) (*Group, error) {
	var group Group
	query := `SELECT id, group_name, st_year from groups WHERE id = $1`
	err := s.client.Get(ctx, &group, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrObjectNotFound
	}

	return &group, err
}

func (s *GroupsRepository) Add(ctx context.Context, group *Group) (int64, error) {
	query := `INSERT INTO groups (group_name, st_year) VALUES($1, $2) RETURNING id`

	var id int64

	err := s.client.ExecQueryRow(ctx, query,
		group.Name.String,
		group.Year).Scan(&id)

	return id, err
}

func (s *GroupsRepository) UpdateById(ctx context.Context, id int64, group *Group) (bool, error) {
	query := `UPDATE groups
			SET group_name = $1, st_year = $2
			WHERE id = $3`

	result, err := s.client.Exec(ctx, query,
		group.Name.String,
		group.Year,
		id)

	return result.RowsAffected() > 0, err
}

func (s *GroupsRepository) Remove(ctx context.Context, id int64) (bool, error) {
	result, err := s.client.Exec(ctx, "DELETE FROM groups WHERE id = $1", id)

	return result.RowsAffected() > 0, err
}
