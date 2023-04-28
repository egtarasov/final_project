package group_client

import (
	"context"
	"google.golang.org/grpc"
	"homework-5/client/pb/group_repo"
)

type StudentClient struct {
	grpc group_repo.GroupServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, target string) (*StudentClient, error) {
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	grpc := group_repo.NewGroupServiceClient(conn)

	return &StudentClient{
		grpc: grpc,
		conn: conn,
	}, nil
}

func (s *StudentClient) Close() error {
	return s.conn.Close()
}
