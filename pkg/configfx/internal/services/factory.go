package services

import (
	"github.com/postmanq/postmanq/pkg/configfx/config"
	gc "go.uber.org/config"
)

func NewFxConfigProviderFactory() config.ProviderFactory {
	return &factory{}
}

type factory struct{}

func (f *factory) Create(options ...config.Option) (config.Provider, error) {
	gcOptions := make([]gc.YAMLOption, len(options))
	for i, option := range options {
		gcOptions[i] = option
	}

	p, err := gc.NewYAML(gcOptions...)
	if err != nil {
		return nil, err
	}

	return &provider{p}, nil
}
