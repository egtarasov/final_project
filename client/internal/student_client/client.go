package student_client

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"homework-5/client/pb/student_repo"
)

type Client interface {
	GetById(ctx context.Context, id int64) (*student_repo.Student, error)
	Create(ctx context.Context, group *student_repo.Student) (int64, error)
	Update(ctx context.Context, id int64, group *student_repo.Student) (bool, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

type StudentClient struct {
	grpc student_repo.StudentServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, target string) (*StudentClient, error) {
	conn, err := grpc.DialContext(ctx,
		target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))

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
