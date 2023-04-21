package group

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupsRepository_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewTestDb()
		db.setUp(t)
		defer db.tearDown(t)

		group := DefaultGroup().P()

		id, err := db.repo.Add(db.ctx, group)

		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
}

func TestGroupsRepository_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewTestDb()
		db.setUp(t)
		defer db.tearDown(t)
		group := DefaultGroup().P()
		var id uint64
		// add test data to db
		db.db.ExecQueryRow(
			db.ctx,
			`INSERT INTO groups (group_name, st_year) VALUES($1, $2) RETURNING id`,
			group.Name,
			group.Year).Scan(&id)
		// update index according to database
		group.Id = int64(id)

		g, err := db.repo.GetById(db.ctx, id)

		require.NoError(t, err)
		assert.Equal(t, group, g)
	})
}

func TestGroupsRepository_Remove(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		db := NewTestDb()
		db.setUp(t)
		defer db.tearDown(t)
		group := DefaultGroup().P()
		var id uint64
		// add test data to db
		db.db.ExecQueryRow(
			db.ctx,
			`INSERT INTO groups (group_name, st_year) VALUES($1, $2) RETURNING id`,
			group.Name,
			group.Year).Scan(&id)

		tt := []struct {
			id     uint64
			expect bool
		}{
			{id, true},
			{id, false},
			{id + 1, false},
			{id + 1000, false},
		}

		for _, tc := range tt {
			ok, err := db.repo.Remove(db.ctx, tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expect, ok)
		}
	})
}

func TestGroupsRepository_UpdateById(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		db := NewTestDb()
		db.setUp(t)
		defer db.tearDown(t)
		group := DefaultGroup()
		var id uint64
		// add test data to db
		db.db.ExecQueryRow(
			db.ctx,
			`INSERT INTO groups (group_name, st_year) VALUES($1, $2) RETURNING id`,
			group.P().Name,
			group.P().Year).Scan(&id)

		tt := []struct {
			id     uint64
			expect bool
			group  *Group
		}{
			{id, true, group.Year(4).P()},
			{id, true, group.Name("NEW_Name").P()},
			{id + 1, false, group.P()},
			{id + 1000, false, group.P()},
		}

		for _, tc := range tt {
			ok, err := db.repo.UpdateById(db.ctx, tc.id, tc.group)
			require.NoError(t, err)
			assert.Equal(t, tc.expect, ok)
		}
	})
}
