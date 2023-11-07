package services

import (
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	gc "go.uber.org/config"
	"os"
)

func NewFxConfigProviderFactory() config.ProviderFactory {
	return &factory{}
}

type factory struct{}

func (f *factory) Create(options ...config.Option) (config.Provider, error) {
	gcOptions := make([]gc.YAMLOption, 0)
	for _, option := range options {
		gcOptions = append(gcOptions, option)
	}

	gcOptions = append(gcOptions, config.Expand(func(key string) (string, bool) {
		return os.LookupEnv(key)
	}))

	p, err := gc.NewYAML(gcOptions...)
	if err != nil {
		return nil, err
	}

	return &provider{p}, nil
}
