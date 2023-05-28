package group_client

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"homework-5/client/pb/group_repo"
)

type Client interface {
	GetById(ctx context.Context, id int64) (*group_repo.Group, error)
	Create(ctx context.Context, group *group_repo.Group) (int64, error)
	Update(ctx context.Context, id int64, group *group_repo.Group) (bool, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

type GroupClient struct {
	grpc group_repo.GroupServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, target string) (*GroupClient, error) {
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		return nil, err
	}
	grpc := group_repo.NewGroupServiceClient(conn)

	return &GroupClient{
		grpc: grpc,
		conn: conn,
	}, nil
}

func (s *GroupClient) Close() error {
	return s.conn.Close()
}
