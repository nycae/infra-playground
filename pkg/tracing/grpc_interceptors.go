package tracing

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

func clientInterceptors(service string) []grpc.DialOption {
	tracer, err := NewOTLP(service)
	if err != nil {
		panic(err)
	}

	options := []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(
			otelgrpc.WithTracerProvider(tracer),
			otelgrpc.WithPropagators(propagation.TraceContext{}))),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(
			otelgrpc.WithTracerProvider(tracer),
			otelgrpc.WithPropagators(propagation.TraceContext{}))),
	}

	return options
}

func serverInterceptors(service string) []grpc.ServerOption {
	tracer, err := NewOTLP(service)
	if err != nil {
		panic(err)
	}

	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithTracerProvider(tracer),
			otelgrpc.WithPropagators(propagation.TraceContext{}))),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(
			otelgrpc.WithTracerProvider(tracer),
			otelgrpc.WithPropagators(propagation.TraceContext{}))),
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
