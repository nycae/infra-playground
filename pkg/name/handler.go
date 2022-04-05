package name

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/nycae/infra-playground/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	url = "https://randomuser.me/api/?results=%d"
	min = 5
	max = 20
)

type Server struct {
	api.UnimplementedNameManagerServer

	names []api.FullName
}

func (s *Server) GetRandom(ctx context.Context, _ *emptypb.Empty) (
	*api.FullName, error) {
	return &s.names[rand.Intn(len(s.names))], nil
}

func (s *Server) GetAll(_ *emptypb.Empty, stream api.NameManager_GetAllServer) error {
	for _, name := range s.names {
		if err := stream.Send(&name); err != nil {
			return fmt.Errorf("error sending name: %v", err)
		}
	}

	return nil
}

func fetchNames() []api.FullName {
	var body struct {
		Results []struct {
			Name struct {
				First string `json:"first"`
				Last  string `json:"last"`
			} `json:"name"`
			Login struct {
				Username string `json:"username"`
			}
		} `json:"results"`
	}

	resp, err := http.Get(fmt.Sprintf(url, rand.Intn(max-min)+min))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		panic(err)
	}

	names := make([]api.FullName, len(body.Results))
	for i, name := range body.Results {
		names[i].FirstName = name.Name.First
		names[i].LastName = name.Login.Username
		names[i].FamilyName = name.Name.Last
	}

	return names
}

func NewServicer() api.NameManagerServer {
	rand.Seed(time.Now().UnixNano())
	return &Server{names: fetchNames()}
}
