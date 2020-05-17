package service

import (
	"github.com/postmanq/postmanq/cli"
	"github.com/postmanq/postmanq/module"
	"go.uber.org/config"
)

func NewConfigProviderByArgs(args cli.Arguments) (module.ConfigProvider, error) {
	return NewConfigProviderByFile(args.ConfigFilename)
}

func NewConfigProviderByOptions(options ...config.YAMLOption) (module.ConfigProvider, error) {
	provider, err := config.NewYAML(options...)
	if err != nil {
		return nil, err
	}

	return &configProvider{provider}, nil
}

func NewConfigProviderByFile(filename string) (module.ConfigProvider, error) {
	return NewConfigProviderByOptions(config.File(filename))
}

type configProvider struct {
	provider config.Provider
}

func (s *configProvider) Populate(key string, target interface{}) error {
	return s.provider.Get(key).Populate(target)
}
