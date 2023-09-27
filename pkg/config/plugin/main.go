package main

import (
	"github.com/postmanq/postmanq/pkg/config/internal/services"
	"github.com/postmanq/postmanq/pkg/plugin"
)

var (
	Plugin plugin.Construct = func() plugin.Plugin {
		return plugin.Plugin{
			Constructs: []interface{}{
				services.NewConfigProviderFactory,
			},
		}
	}
)
