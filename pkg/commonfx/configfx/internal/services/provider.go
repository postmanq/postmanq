package services

import (
	gc "go.uber.org/config"
)

type provider struct {
	provider gc.Provider
}

func (s *provider) Populate(target interface{}) error {
	return s.PopulateByKey("", target)
}

func (s *provider) PopulateByKey(key string, target interface{}) error {
	return s.provider.Get(key).Populate(target)
}
