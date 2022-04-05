package height

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/nycae/infra-playground/api"
	"github.com/nycae/infra-playground/pkg/tracing"
	"github.com/nycae/infra-playground/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	cmInAnInch = 2.54
)

var (
	limiterHost = utils.GetEnvWithDefault("LIMITER_HOST", "limiter:8086")
)

type Server struct {
	api.UnimplementedHeightGuesserServer
	limiter api.HeightLimiterClient
}

func (s *Server) GuessHeight(stream api.HeightGuesser_GuessHeightServer) error {
	limits, err := s.limiter.GetHeightLimits(stream.Context(), &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("unable to comunicate with limiter server: %v", err.Error())
	}

	for _, err := stream.Recv(); err != io.EOF; _, err = stream.Recv() {
		h := rand.Int31n(limits.HeightMax-limits.HeightMin) + limits.HeightMin
		height := api.HeightReport{
			HeightCm: h,
			HeightIn: int32(float32(h) / cmInAnInch),
			HeightFt: 0,
		}

		for height.HeightIn > 12 {
			height.HeightFt += 1
			height.HeightIn -= 12
		}

		stream.Send(&height)
	}

	return nil
}

func NewServicer() api.HeightGuesserServer {
	rand.Seed(time.Now().UnixNano())
	dial, err := grpc.Dial(limiterHost, append(tracing.LimiterClientInterceptors(),
		grpc.WithInsecure())...)
	if err != nil {
		log.Fatal(err)
	}
	return &Server{limiter: api.NewHeightLimiterClient(dial)}
}
