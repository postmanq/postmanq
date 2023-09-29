package main

import (
	"github.com/postmanq/postmanq/pkg/config"
	"github.com/postmanq/postmanq/pkg/config/internal/services"
	"github.com/postmanq/postmanq/pkg/plugin"
)

var (
	Plugin plugin.Construct = func() plugin.Plugin {
		return plugin.Plugin{
			Constructs: []interface{}{
				services.NewFxConfigProviderFactory,
				NewFxConfigProvider,
				NewFxConfig,
			},
		}
	}
)

func NewFxConfigProvider(
	factory config.ProviderFactory,
	filename string,
) (config.Provider, error) {
	return factory.Create(config.File(filename))
}

func NewFxConfig(
	provider config.Provider,
) (*config.Config, error) {
	cfg := new(config.Config)
	err := provider.Populate(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
