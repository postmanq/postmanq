package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/rabbitmq/component"
	"github.com/postmanq/postmanq/module/rabbitmq/service"
)

var (
	PqModule module.DescriptorConstruct = func() module.Descriptor {
		return module.Descriptor{
			Constructs: []interface{}{
				service.NewPool,
				component.NewReceiver,
			},
		}
	}
)
