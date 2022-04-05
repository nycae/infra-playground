package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nycae/infra-playground/pkg/frontend"
	"github.com/nycae/infra-playground/pkg/tracing"
	"github.com/nycae/infra-playground/pkg/utils"
)

const (
	host = "0.0.0.0"
	port = 8080
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("/static")))
	mux.Handle("/v1/", tracing.HandlerMiddleware("frontend-server",
		"home-page", frontend.NewHandler()))

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	log.Printf("frontend server is up and running on port %d", port)
	utils.WaitForShutdown()

	ctx, cancelContext := context.WithTimeout(context.Background(), time.Second*20)
	if err := srv.Shutdown(ctx); err != nil {
		cancelContext()
		log.Fatal(err)
	}
}
