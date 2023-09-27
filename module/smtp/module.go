package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/smtp/component"
	"github.com/postmanq/postmanq/module/smtp/service"
)

var (
	PqModule module.PluginConstruct = func() module.Plugin {
		return module.Plugin{
			Constructs: []interface{}{
				service.NewConnectorFactory,
				service.NewScanner,
				component.NewSender,
			},
		}
	}
)
