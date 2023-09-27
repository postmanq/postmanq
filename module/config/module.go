package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/config/service"
)

var (
	PqModule module.PluginConstruct = func() module.Plugin {
		return module.Plugin{
			Constructs: []interface{}{
				service.NewConfigProviderFactory,
				service.NewConfigProviderByArgs,
			},
		}
	}
)
