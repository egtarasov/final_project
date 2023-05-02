package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"homework-5/server/internal/app"
	"homework-5/server/internal/app/database"
	"homework-5/server/internal/app/group"
	"homework-5/server/internal/app/group_service"
	"homework-5/server/internal/app/pb/group_repo"
	"homework-5/server/internal/app/pb/student_repo"
	"homework-5/server/internal/app/student"
	"homework-5/server/internal/app/student_service"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	service     = "repo_service"
	environment = "development"
)

func Tracer(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tr := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
		)),
	)

	return tr, nil
}

const (
	jaegerUrl = "http://localhost:14268/api/traces"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tr, err := Tracer(jaegerUrl)
	if err != nil {
		log.Fatalln(err)
	}

	otel.SetTracerProvider(tr)

	// Graceful shutdown
	defer func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tr.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}(ctx)

	// Start prometheus
	app.Init()
	go http.ListenAndServe(":9091", promhttp.Handler())

	db, err := database.NewDb(ctx)
	if err != nil {
		log.Fatalln("cant connect to database")
	}

	studentRepo := student.NewStudentsRepository(db)
	groupRepo := group.NewGroupsRepository(db)

	// Start student service
	go func() {
		lsn, err := net.Listen("tcp", ":8000")
		if err != nil {
			log.Fatalln(err)
		}
		server := grpc.NewServer()
		implementation := student_service.NewImplementation(studentRepo)
		student_repo.RegisterStudentServiceServer(server, implementation)

		if err = server.Serve(lsn); err != nil {
			log.Fatal(err)
		}
	}()

	// Start group server
	lsn, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))

	implementation := group_service.NewGroupService(groupRepo)
	group_repo.RegisterGroupServiceServer(server, implementation)

	if err = server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
