package main

import (
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/queue/component"
	qs "github.com/postmanq/postmanq/module/queue/service"
	"go.uber.org/fx"
)

type PqModuleIn struct {
	fx.In
	ConfigProvider cs.ConfigProvider
}

type PqModuleOut struct {
	fx.Out
	Pool     qs.Pool
	Receiver interface{} `group:"component"`
}

func PqModule(params PqModuleIn) PqModuleOut {
	pool := qs.NewPool()
	return PqModuleOut{
		Pool: pool,
		Receiver: component.NewReceiver(
			params.ConfigProvider,
			pool,
		),
	}
}
