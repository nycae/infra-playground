package age

import (
	"context"
	"fmt"
	"sync"

	"github.com/nycae/infra-playground/api"
)

type ageCache struct {
	lock sync.RWMutex
	ages map[string]int32
}

func (a *ageCache) getAge(key string) (int32, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	res, ok := a.ages[key]

	return res, ok
}

func (a *ageCache) setAge(key string, val int32) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.ages[key] = val
}

type Server struct {
	api.UnimplementedAgeManagerServer

	api   *API
	cache ageCache
}

func (s *Server) GetBirthdayOf(ctx context.Context, name *api.FullName) (
	*api.AgeReport, error) {
	age, ok := s.cache.getAge(name.FirstName)
	if !ok {
		fetchedAge, err := s.api.AgeOf(name.FirstName)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch the age of user: %v", err.Error())
		}

		age = fetchedAge
		s.cache.setAge(name.FamilyName, age)
	}

	return &api.AgeReport{Age: age, Birthday: nil}, nil
}

func NewServicer() api.AgeManagerServer {
	return &Server{
		api:   NewAPI(),
		cache: ageCache{ages: make(map[string]int32)},
	}
}
