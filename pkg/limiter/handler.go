package limiter

import (
	"context"

	"github.com/nycae/infra-playground/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	api.UnimplementedHeightLimiterServer
}

const (
	limitMax = 220
	limitMin = 140
)

func (s *Server) GetHeightLimits(ctx context.Context, _ *emptypb.Empty) (
	*api.HeightLimits, error) {
	return &api.HeightLimits{HeightMax: limitMax, HeightMin: limitMin}, nil
}

func NewServicer() api.HeightLimiterServer {
	return &Server{}
}
