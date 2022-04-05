package age

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultEndpoint = "https://api.agify.io/"
)

type API struct {
	host   *url.URL
	client *http.Client
}

func (api *API) AgeOf(name string) (int32, error) {
	api.host.Query().Add("name", name)
	defer func() {
		for key := range api.host.Query() {
			api.host.Query().Del(key)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, api.host.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("unable to build request: %v", err.Error())
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("invalid request: %v", err.Error())
	}
	defer resp.Body.Close()

	var body struct {
		Name  string `json:"name"`
		Age   int32  `json:"age"`
		Count int32  `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("unable to parse api response: %v", err.Error())
	}

	return body.Age, nil
}

func NewAPI() *API {
	host, _ := url.Parse(defaultEndpoint)
	return &API{
		host:   host,
		client: &http.Client{},
	}
}
