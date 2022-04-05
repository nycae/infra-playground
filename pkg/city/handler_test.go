package city

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/nycae/infra-playground/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	addr = "0.0.0.0"
	port = 50055
)

func server() *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	api.RegisterCityManagerServer(srv, NewServicer())

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Print(err)
		}
	}()

	return srv
}

func Test_API(t *testing.T) {
	srv := server()
	defer srv.GracefulStop()

	dial, _ := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithInsecure())
	servicer := api.NewCityManagerClient(dial)

	t.Run("default", func(t *testing.T) {
		if _, err := servicer.GetRandomCity(context.Background(),
			&emptypb.Empty{}); err != nil {
			t.Error(err)
		}

		city, err := servicer.GetCityAt(context.Background(), &api.Location{
			Latitude: 38.9848, Longitude: 3.9274})
		if err != nil {
			t.Error(err)
		}

		if expected := "Ciudad Real"; city.Name != expected {
			t.Errorf("expected name %v, have %v", expected, city.Name)
		}
	})
}
