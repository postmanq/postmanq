package services

import (
	"github.com/postmanq/postmanq/pkg/config"
)

func NewConfigProviderFactory() config.ProviderFactory {
	return &factory{}
}

type factory struct{}

func (f *factory) CreateFromMap(data map[string]string) (config.Provider, error) {
	return NewConfigProviderByOptions(config.Static(data))
}
