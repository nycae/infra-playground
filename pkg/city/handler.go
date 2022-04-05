package city

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/nycae/infra-playground/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

var cities = []*api.City{
	{
		Name:     "Ciudad Real",
		Location: &api.Location{Latitude: 38.9848, Longitude: 3.9274},
	},
	{
		Name:     "Nürnberg",
		Location: &api.Location{Latitude: 49.4521, Longitude: 11.0767},
	},
}

type Server struct {
	api.UnimplementedCityManagerServer

	cities []*api.City
}

func (s *Server) GetCityAt(ctx context.Context, loc *api.Location) (*api.City, error) {
	for _, city := range s.cities {
		if city.Location.Longitude == loc.Longitude &&
			city.Location.Latitude == loc.Latitude {
			return city, nil
		}
	}

	return nil, fmt.Errorf("no city at [%v, %v]", loc.Latitude, loc.Longitude)
}

func (s *Server) GetRandomCity(ctx context.Context, _ *emptypb.Empty) (
	*api.City, error) {
	return s.cities[rand.Intn(len(s.cities))], nil
}

func NewServicer() api.CityManagerServer {
	return &Server{cities: cities}
}
