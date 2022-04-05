package tracing

import (
	"fmt"

	"github.com/nycae/infra-playground/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type url struct {
	host, port string
}

func (u url) String() string {
	return fmt.Sprintf("http://%s:%s/api/traces", u.host, u.port)
}

var (
	defaultURL = url{
		host: utils.GetEnvWithDefault("JAEGER_HOST", "jaeger"),
		port: utils.GetEnvWithDefault("JAEGER_PORT", "14268"),
	}
)

func clientInterceptors(service string) []grpc.DialOption {
	j, err := NewJaegerExporter()
	if err != nil {
		panic(err)
	}

	tracer := NewTraceProvider(service, j)
	options := []grpc.DialOption{
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(
			otelgrpc.WithTracerProvider(tracer),
		)),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(
			otelgrpc.WithTracerProvider(tracer),
		)),
	}

	return options
}

func serverInterceptors(service string) []grpc.ServerOption {
	j, err := NewJaegerExporter()
	if err != nil {
		panic(err)
	}

	tracer := NewTraceProvider(service, j)

	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithTracerProvider(tracer),
		)),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(
			otelgrpc.WithTracerProvider(tracer),
		)),
	}

	return options
}

func AgerClientInterceptors() []grpc.DialOption {
	return clientInterceptors("ager-client")
}

func NamerClientInterceptors() []grpc.DialOption {
	return clientInterceptors("namer-client")
}

func CitierClientInterceptors() []grpc.DialOption {
	return clientInterceptors("citier-client")
}

func HeighterClientInterceptors() []grpc.DialOption {
	return clientInterceptors("heighter-client")
}

func LimiterClientInterceptors() []grpc.DialOption {
	return clientInterceptors("limiter-client")
}

func AgerServerInterceptors() []grpc.ServerOption {
	return serverInterceptors("ager-server")
}

func NamerServerInterceptors() []grpc.ServerOption {
	return serverInterceptors("namer-server")
}

func CitierServerInterceptors() []grpc.ServerOption {
	return serverInterceptors("citier-server")
}

func HeighterServerInterceptors() []grpc.ServerOption {
	return serverInterceptors("heighter-server")
}

func LimiterServerInterceptors() []grpc.ServerOption {
	return serverInterceptors("limiter-server")
}
