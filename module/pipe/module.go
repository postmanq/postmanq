package main

import (
	"github.com/postmanq/postmanq/module/pipe/component"
	"go.uber.org/fx"
)

type PqModuleIn struct {
	fx.In
	Components []interface{} `group:"component"`
}

type PqModuleOut struct {
	fx.Out
	Pipe *component.Pipe
}

func PqModule(params PqModuleIn) PqModuleOut {
	return PqModuleOut{
		Pipe: component.NewPipe(params.Components),
	}
}
