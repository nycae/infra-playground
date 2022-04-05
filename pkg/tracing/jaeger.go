package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewJaegerExporter() (trace.SpanExporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(defaultURL.String())))
	if err != nil {
		return nil, fmt.Errorf("unable to contact jaeger collector: %v", err.Error())
	}
	return exp, nil
}
