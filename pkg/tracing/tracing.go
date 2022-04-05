package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func NewTraceProvider(server string, exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	tr := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(server),
		)),
	)

	otel.SetTracerProvider(tr)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tr
}
