package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/component"
	"github.com/postmanq/postmanq/module/pipe/service"
)

var (
	PqModule module.DescriptorConstruct = func() module.Descriptor {
		return module.Descriptor{
			Constructs: []interface{}{
				component.NewRunner,
				service.NewComponentFactory,
				service.NewStageFactory,
				service.NewPipelineFactory,
			},
		}
	}
)
