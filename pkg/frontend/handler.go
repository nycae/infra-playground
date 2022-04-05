package frontend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sync"

	"github.com/nycae/infra-playground/api"
	"github.com/nycae/infra-playground/pkg/tracing"
	"github.com/nycae/infra-playground/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	indexTemplate = `
<head>
  <link href="https://unpkg.com/material-components-web@latest/dist/material-components-web.min.css" rel="stylesheet">
  <script src="https://unpkg.com/material-components-web@latest/dist/material-components-web.min.js"></script>
  <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
</head>
<body>
  <header lass="mdc-top-app-bar">
    <div class="mdc-top-app-bar__row">
      <section class="mdc-top-app-bar__section mdc-top-app-bar__section--align-start">
        <span class="mdc-top-app-bar__title">{{ .Message }}</span>
      </section>
    </div>
	</header>
	<main class="mdc-top-app-bar--fixed-adjust">
		<ul class="mdc-list">
	  {{ range $person := .People }}
		  <li class="mdc-list-item" tabindex="0">
				<span class="mdc-list-item__text">
          <span class="mdc-list-item__primary-text">
					  Person: 
            {{ $person.Name.FamilyName }},
					  {{ $person.Name.FirstName }}
			    {{ if $person.Name.LastName }}
            {{ $person.Name.LastName }}
			    {{ end }}
          </span>
          <span class="mdc-list-item__secondary-text">
						Lives in: {{ $person.Residence.Name }}	
						Height: {{ $person.Height }}
          </span>
				</span>
			</li>
			<li role="separator" class="mdc-list-divider"></li>
	  {{ end }}
		</ul>
  </main>
</body>
`
)

var (
	nameEndpoint   = utils.GetEnvWithDefault("NAME_SERVICE_HOST", "localhost:8085")
	cityEndpoint   = utils.GetEnvWithDefault("CITY_SERVICE_HOST", "localhost:8082")
	ageEndpoint    = utils.GetEnvWithDefault("AGE_SERVICE_HOST", "localhost:8083")
	heightEndpoint = utils.GetEnvWithDefault("HEIGHT_SERVICE_HOST", "localhost:8084")
)

type Handler struct {
	index *template.Template

	nameService   api.NameManagerClient
	cityService   api.CityManagerClient
	ageService    api.AgeManagerClient
	heightService api.HeightGuesserClient
}

func (h *Handler) sendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}

func (h *Handler) fetchPeopleNames(ctx context.Context) ([]*api.FullName, error) {
	names, err := h.nameService.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch names: %v", err.Error())
	}

	var nn []*api.FullName
	for name, err := names.Recv(); err != io.EOF; name, err = names.Recv() {
		if err != nil {
			return nil, fmt.Errorf("error fetching names from server: %v", err.Error())
		}
		nn = append(nn, name)
	}

	return nn, nil
}

func (h *Handler) fetchPeopleHeigts(ctx context.Context, names []*api.FullName) (
	[]*api.HeightReport, error) {
	stream, err := h.heightService.GuessHeight(ctx)
	defer stream.CloseSend()
	if err != nil {
		return nil, fmt.Errorf("unable to create stream: %v", err.Error())
	}

	hh := make([]*api.HeightReport, len(names))
	for i, name := range names {
		if err := stream.Send(name); err != nil {
			return nil, fmt.Errorf("error sending message: %v", err.Error())
		}

		hh[i], err = stream.Recv()
		if err != nil {
			return nil, fmt.Errorf("error receiving message: %v", err.Error())
		}
	}

	return hh, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nn, err := h.fetchPeopleNames(ctx)
	if err != nil {
		h.sendError(w, fmt.Errorf("error fetching names: %v", err.Error()))
		return
	}

	hh, err := h.fetchPeopleHeigts(ctx, nn)
	if err != nil {
		h.sendError(w, fmt.Errorf("error fetching heights: %v", err.Error()))
		return
	}

	if len(nn) != len(hh) {
		h.sendError(w, errors.New("unable to fetch all names or heights"))
		return
	}

	var wait sync.WaitGroup
	wait.Add(len(nn))

	pp := make([]api.Person, len(nn))
	for i, n := range nn {
		go func(ctx context.Context, n *api.FullName, i int) {
			defer wait.Done()
			age, err := h.ageService.GetBirthdayOf(ctx, n)
			if err != nil {
				h.sendError(w, fmt.Errorf("unable to fetch age: %v", err.Error()))
				return
			}

			pp[i].Name = n
			pp[i].Height = hh[i].HeightCm
			pp[i].Birthday = age.Birthday
			pp[i].Residence, err = h.cityService.GetRandomCity(ctx, &emptypb.Empty{})
			if err != nil {
				h.sendError(w, fmt.Errorf("unable to fetch city: %v", err.Error()))
				return
			}
		}(ctx, n, i)
	}

	wait.Wait()

	h.index.ExecuteTemplate(w, "index", map[string]interface{}{
		"People":  pp,
		"Message": "Welcome to the people list example",
	})
}

func NewHandler() http.Handler {
	tmpl := template.Must(template.New("index").Parse(indexTemplate))

	nameManager := api.NewNameManagerClient(utils.MustGetDial(nameEndpoint,
		append(tracing.NamerClientInterceptors(), grpc.WithInsecure())...))
	cityManager := api.NewCityManagerClient(utils.MustGetDial(cityEndpoint,
		append(tracing.CitierClientInterceptors(), grpc.WithInsecure())...))
	ageManager := api.NewAgeManagerClient(utils.MustGetDial(ageEndpoint,
		append(tracing.AgerClientInterceptors(), grpc.WithInsecure())...))
	heightManager := api.NewHeightGuesserClient(utils.MustGetDial(heightEndpoint,
		append(tracing.HeighterClientInterceptors(), grpc.WithInsecure())...))

	return &Handler{
		index:         tmpl,
		nameService:   nameManager,
		cityService:   cityManager,
		ageService:    ageManager,
		heightService: heightManager,
	}
}
