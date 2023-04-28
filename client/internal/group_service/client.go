package student_client

import (
	"context"
	"google.golang.org/grpc"
	"homework-5/client/pb/student_repo"
)

type StudentClient struct {
	grpc student_repo.StudentServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, target string) (*StudentClient, error) {
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	grpc := student_repo.NewStudentServiceClient(conn)

	return &StudentClient{
		grpc: grpc,
		conn: conn,
	}, nil
}

func (s *StudentClient) Close() error {
	return s.conn.Close()
}
