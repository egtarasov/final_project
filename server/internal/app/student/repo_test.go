package student

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

		student := DefaultStudent().P()
		var id uint64
		db.add(student, &id)
		id, err := db.repo.Add(db.ctx, student)

		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
}

func TestGroupsRepository_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewTestDb()
		db.setUp(t)
		defer db.tearDown(t)
		student := DefaultStudent().P()
		var id uint64
		// add test data to db
		db.add(student, &id)

		// update index according to database
		student.Id = int64(id)

		g, err := db.repo.GetById(db.ctx, id)

		require.NoError(t, err)
		assert.Equal(t, student, g)
	})
}

func TestGroupsRepository_Remove(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		db := NewTestDb()
		db.setUp(t)
		defer db.tearDown(t)
		student := DefaultStudent().P()
		var id uint64
		// add test data to db
		db.add(student, &id)

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
		student := DefaultStudent()
		var id uint64
		// add test data to db
		db.add(student.P(), &id)
		student.Id(int64(id))
		tt := []struct {
			id     uint64
			expect bool
			group  *Student
		}{
			{id, true, student.Gpa(4).P()},
			{id, true, student.FirstName("Oleg").P()},
			{id + 1, false, student.P()},
			{id + 1000, false, student.P()},
		}

		for _, tc := range tt {
			ok, err := db.repo.UpdateById(db.ctx, tc.id, tc.group)
			require.NoError(t, err)
			assert.Equal(t, tc.expect, ok)
		}
	})
}
