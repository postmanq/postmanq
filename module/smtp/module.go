package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/smtp/component"
	"github.com/postmanq/postmanq/module/smtp/service"
)

var (
	PqModule module.DescriptorConstruct = func() module.Descriptor {
		return module.Descriptor{
			Constructs: []interface{}{
				service.NewConnectorFactory,
				service.NewScanner,
				component.NewSender,
			},
		}
	}
)
