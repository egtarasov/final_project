package server

import (
	"context"
	"github.com/golang/mock/gomock"
	mock "homework-5/internal/app/server/mocks"
	"testing"
)

type UserRepoFixture struct {
	ctx         context.Context
	ctrl        *gomock.Controller
	s           *server
	studentRepo *mock.MockStudentsRepository
	groupRepo   *mock.MockGroupsRepository
}

func setUp(t *testing.T) *UserRepoFixture {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	studentRepo := mock.NewMockStudentsRepository(ctrl)
	groupRepo := mock.NewMockGroupsRepository(ctrl)
	s := NewServer(ctx, studentRepo, groupRepo)
	return &UserRepoFixture{
		ctx:         ctx,
		ctrl:        ctrl,
		s:           s,
		studentRepo: studentRepo,
		groupRepo:   groupRepo,
	}
}

func (u *UserRepoFixture) teraDown() {
	u.ctx.Done()
	u.ctrl.Finish()
}
