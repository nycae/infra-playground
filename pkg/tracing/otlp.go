package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"google.golang.org/grpc"
)

func NewOTLP(service string) (*trace.TracerProvider, error) {
	ctx := context.Background()

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String(service)))
	if err != nil {
		return nil, err
	}

	conn, err := grpc.DialContext(ctx, defaultURL.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	te, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := trace.NewBatchSpanProcessor(te)
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(bsp),
		trace.WithResource(res),
	)

	return tp, nil
}
