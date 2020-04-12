package main

import (
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/pipe/component"
	"go.uber.org/fx"
)

type PqModuleIn struct {
	fx.In
	ConfigProvider cs.ConfigProvider
	Components     []interface{} `group:"component"`
}

type PqModuleOut struct {
	fx.Out
	Runner *component.Runner
}

func PqModule(params PqModuleIn) PqModuleOut {
	return PqModuleOut{
		Runner: component.NewRunner(params.ConfigProvider, params.Components),
	}
}
