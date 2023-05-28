package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"homework-5/client/internal/client"
	group_client "homework-5/client/internal/group_client"
	"homework-5/client/internal/student_client"
	"log"
	"net/http"
	"time"
)

const (
	serverStudentUrl = ":8000"
	serverGroupUrl   = ":8001"
)

const (
	environment = "Development"
	name        = "client_repo"
)

func Tracer(target string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(target)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
			attribute.String("environment", environment),
		)),
	)

	return tp, nil
}

const (
	jaegerUrl = "http://localhost:14268/api/traces"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create Tracer Provider
	tp, err := Tracer(jaegerUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Graceful shutdown
	defer func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}(ctx)

	//Register Tracer Provider
	otel.SetTracerProvider(tp)

	studentClient, err := student_client.NewClient(ctx, serverStudentUrl)
	if err != nil {
		log.Fatal("cant connect to student service")
	}
	defer studentClient.Close()

	groupClient, err := group_client.NewClient(ctx, serverGroupUrl)
	if err != nil {
		log.Fatal("cant connect to student service")
	}
	defer groupClient.Close()

	cl := client.NewClient(ctx, groupClient, studentClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/client", cl.Handle)
	if err := http.ListenAndServe(":1200", mux); err != nil {
		log.Fatal(err)
	}
}
