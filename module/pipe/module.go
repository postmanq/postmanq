package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/component"
	"github.com/postmanq/postmanq/module/pipe/service"
	"github.com/postmanq/postmanq/module/pipe/service/factory"
	"github.com/postmanq/postmanq/module/pipe/service/stage"
)

var (
	PqModule module.PluginConstruct = func() module.Plugin {
		return module.Plugin{
			Constructs: []interface{}{
				component.NewRunner,
				factory.NewComponent,
				factory.NewStage,
				service.NewPipelineFactory,
				stage.NewComplete,
				stage.NewMiddleware,
				stage.NewParallelMiddleware,
				stage.NewReceive,
			},
		}
	}
)
