package group

import "context"

type Repository interface {
	Add(ctx context.Context, group *Group) (int64, error)                 //create
	GetById(ctx context.Context, id int64) (*Group, error)                //read
	UpdateById(ctx context.Context, id int64, group *Group) (bool, error) // update
	Remove(ctx context.Context, id int64) (bool, error)                   //delete
}
