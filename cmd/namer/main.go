package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nycae/infra-playground/api"
	"github.com/nycae/infra-playground/pkg/name"
	"github.com/nycae/infra-playground/pkg/tracing"
	"github.com/nycae/infra-playground/pkg/utils"
	"google.golang.org/grpc"
)

const (
	addr = "0.0.0.0"
	port = 8085
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer(tracing.NamerServerInterceptors()...)
	api.RegisterNameManagerServer(srv, name.NewServicer())

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Print(err)
		}
	}()

	log.Printf("Server is listening on port %d", port)
	utils.WaitForShutdown()
	srv.GracefulStop()
}
