package service

import (
	"go.uber.org/config"
)

type ConfigProvider interface {
	Populate(string, interface{}) error
}

func NewConfigProviderByOptions(options ...config.YAMLOption) (ConfigProvider, error) {
	provider, err := config.NewYAML(options...)
	if err != nil {
		return nil, err
	}

	return &configProvider{provider}, nil
}

type configProvider struct {
	provider config.Provider
}

func (s *configProvider) Populate(key string, target interface{}) error {
	return s.provider.Get(key).Populate(target)
}
