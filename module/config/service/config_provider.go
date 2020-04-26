package service

import (
	"github.com/postmanq/postmanq/cli"
	"go.uber.org/config"
)

type ConfigProvider interface {
	Populate(string, interface{}) error
}

func NewConfigProviderByArgs(args cli.Arguments) (ConfigProvider, error) {
	return NewConfigProviderByFile(args.ConfigFilename)
}

func NewConfigProviderByOptions(options ...config.YAMLOption) (ConfigProvider, error) {
	provider, err := config.NewYAML(options...)
	if err != nil {
		return nil, err
	}

	return &configProvider{provider}, nil
}

func NewConfigProviderByFile(filename string) (ConfigProvider, error) {
	return NewConfigProviderByOptions(config.File(filename))
}

type configProvider struct {
	provider config.Provider
}

func (s *configProvider) Populate(key string, target interface{}) error {
	return s.provider.Get(key).Populate(target)
}
