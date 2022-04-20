package tracing

import (
	"fmt"

	"github.com/nycae/infra-playground/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type url struct {
	host, port string
}

func (u url) String() string {
	return fmt.Sprintf("%s:%s", u.host, u.port)
}

var (
	defaultURL = url{
		host: utils.GetEnvWithDefault("COLLECTOR_HOST", "jaeger"),
		port: utils.GetEnvWithDefault("COLLECTOR_PORT", "14268"),
	}
)

func NewTraceProvider(server string,
	exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
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
