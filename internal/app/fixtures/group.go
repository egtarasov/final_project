package fixtures

import (
	"database/sql"
	"homework-5/internal/app/group"
)

type GroupBuilder struct {
	group group.Group
}

func (g *GroupBuilder) Name(name string) *GroupBuilder {
	g.group.Name = sql.NullString{String: name, Valid: true}
	return g
}

func (g *GroupBuilder) Id(id int64) *GroupBuilder {
	g.group.Id = id
	return g
}

func (g *GroupBuilder) Year(year int32) *GroupBuilder {
	g.group.Year = year
	return g
}

func (g *GroupBuilder) V() group.Group {
	return g.group
}

func (g *GroupBuilder) P() *group.Group {
	return &g.group
}

func DefaultGroup() *GroupBuilder {
	group := GroupBuilder{}
	return group.
		Name("BSE-229").
		Id(1).
		Year(3)
}
