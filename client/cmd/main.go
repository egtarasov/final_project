package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	group_client "homework-5/client/internal/group_service"
	"homework-5/client/internal/student_client"
	"log"
	"time"
)

const (
	serverStudentUrl = ":8000"
	serverGroupUrl   = ":80001"
)

const (
	environment = "Development"
	name        = "Client_repo"
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

	studentCLient, err := student_client.NewClient(ctx, serverStudentUrl)
	if err != nil {
		log.Fatal("cant connect to student service")
	}
	defer studentCLient.Close()

	groupClient, err := group_client.NewClient(ctx, serverStudentUrl)
	if err != nil {
		log.Fatal("cant connect to student service")
	}

	defer groupClient.Close()

	//student, err := studentCLient.GetById(ctx, 2)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(student)
	//
	//student.FirstName = "Bob"
	//student.SecondName = "Bobuk"
	//id, err := studentCLient.Create(ctx, student)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(id)
	//
	//ok, err := studentCLient.Delete(ctx, 1)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(ok)
	//
	//time.Sleep(time.Minute * 10)
}
