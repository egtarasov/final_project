package group

import (
	"context"
	"fmt"
	"homework-5/internal/app/database"
	"strings"
	"sync"
	"testing"
)

type TestDb struct {
	sync.Mutex
	ctx  context.Context
	db   *database.Database
	repo *GroupsRepository
}

func NewTestDb() *TestDb {
	ctx := context.Background()
	db, err := database.NewDb(ctx)
	if err != nil {
		return nil
	}
	return &TestDb{
		ctx:  ctx,
		repo: NewGroupsRepository(db),
		db:   db,
	}
}

func (d *TestDb) setUp(t *testing.T) {
	d.Lock()
	d.truncate()
}

func (d *TestDb) tearDown(t *testing.T) {
	defer d.Unlock()
	d.truncate()
}

func (d *TestDb) truncate() {
	var tables []string
	err := d.db.Select(d.ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'goose_db_version'")
	if err != nil {
		return
	}
	if len(tables) == 0 {
		panic("not tables in db")
	}
	query := fmt.Sprintf("TRUNCATE TABLE %v", strings.Join(tables, ", "))
	if _, err := d.db.Exec(d.ctx, query); err != nil {
		panic(err)
	}
}
