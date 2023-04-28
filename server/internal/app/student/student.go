package student

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, student *Student) (int64, error)                 //create
	GetById(ctx context.Context, id int64) (*Student, error)                  //read
	UpdateById(ctx context.Context, id int64, student *Student) (bool, error) // update
	Remove(ctx context.Context, id int64) (bool, error)                       //delete
}
