package tracing

import (
	"log"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type handler struct {
	t trace.Tracer
	h http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.t.Start(r.Context(), r.RequestURI,
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	h.h.ServeHTTP(w, r.WithContext(ctx))
}

func HandlerMiddleware(server, service string, h http.Handler) http.Handler {
	j, err := NewOTLP(service)
	if err != nil {
		log.Printf("unable to create tracer: %v", err.Error())
	}

	return &handler{
		h: h,
		t: j.Tracer(service),
	}
}
