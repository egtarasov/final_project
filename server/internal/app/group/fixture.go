package group

import (
	"database/sql"
)

type GroupBuilder struct {
	group Group
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

func (g *GroupBuilder) V() Group {
	return g.group
}

func (g *GroupBuilder) P() *Group {
	return &g.group
}

func DefaultGroup() *GroupBuilder {
	group := GroupBuilder{}
	return group.
		Name("BSE-229").
		Id(1).
		Year(3)
}
