package service

import (
	"github.com/postmanq/postmanq/module"
	"go.uber.org/config"
)

type ConfigProviderFactory interface {
	CreateFromMap(map[string]string) (module.ConfigProvider, error)
}

func NewConfigProviderFactory() ConfigProviderFactory {
	return &factory{}
}

type factory struct{}

func (f *factory) CreateFromMap(data map[string]string) (module.ConfigProvider, error) {
	return NewConfigProviderByOptions(config.Static(data))
}
