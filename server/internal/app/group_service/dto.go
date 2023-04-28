package group_service

import (
	"database/sql"
	"homework-5/internal/app/group"
	"homework-5/internal/app/pb/group_repo"
)

func ParseGroup(gr *group_repo.Group) *group.Group {
	return &group.Group{
		Id:   gr.Id,
		Name: sql.NullString{String: gr.Name, Valid: true},
		Year: gr.Year,
	}
}

func ParseGroupRequest(gr *group.Group) *group_repo.Group {
	return &group_repo.Group{
		Id:   gr.Id,
		Name: gr.Name.String,
		Year: gr.Year,
	}
}
