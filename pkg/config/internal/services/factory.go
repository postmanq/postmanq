package services

import (
	"github.com/postmanq/postmanq/pkg/config"
)

func NewFxConfigProviderFactory() config.ProviderFactory {
	return &factory{}
}

type factory struct{}

func (f *factory) Create(options ...config.Option) (config.Provider, error) {
	return NewConfigProviderByOptions(options...)
}
