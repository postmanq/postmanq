package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/component"
	"github.com/postmanq/postmanq/module/pipe/service/factory"
	"github.com/postmanq/postmanq/module/pipe/service/stage"
)

var (
	PqModule module.DescriptorConstruct = func() module.Descriptor {
		return module.Descriptor{
			Constructs: []interface{}{
				component.NewRunner,
				factory.NewComponentFactory,
				factory.NewStageFactory,
				factory.NewPipelineFactory,
				stage.NewComplete,
				stage.NewMiddleware,
				stage.NewParallelMiddleware,
				stage.NewReceive,
			},
		}
	}
)
