package services

import (
	"github.com/postmanq/postmanq/pkg/config"
	gc "go.uber.org/config"
)

func NewConfigProviderByOptions(options ...config.Option) (config.Provider, error) {
	gcOptions := make([]gc.YAMLOption, len(options))
	for i, option := range options {
		gcOptions[i] = option
	}

	provider, err := gc.NewYAML(gcOptions...)
	if err != nil {
		return nil, err
	}

	return &configProvider{provider}, nil
}

type configProvider struct {
	provider gc.Provider
}

func (s *configProvider) Populate(target interface{}) error {
	return s.PopulateByKey("", target)
}

func (s *configProvider) PopulateByKey(key string, target interface{}) error {
	return s.provider.Get(key).Populate(target)
}
