package main

import (
	"github.com/postmanq/postmanq/module"
)

var (
	PqModule module.DescriptorConstruct = func() module.Descriptor {
		return module.Descriptor{
			Constructs: []interface{}{},
		}
	}
)
