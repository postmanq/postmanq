package main

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/validator/service"
)

var (
	PqModule module.DescriptorConstruct = func() module.Descriptor {
		return module.Descriptor{
			Constructs: []interface{}{
				service.NewValidator,
			},
		}
	}
)
