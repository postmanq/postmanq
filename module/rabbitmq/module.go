package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/rabbitmq/component"
	"github.com/postmanq/postmanq/module/rabbitmq/service"
)

var (
	PqModule module.PluginConstruct = func() module.Plugin {
		return module.Plugin{
			Constructs: []interface{}{
				service.NewPool,
				component.NewReceiver,
			},
		}
	}
)
