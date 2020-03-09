package main

import (
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/queue/component"
	qs "github.com/postmanq/postmanq/module/queue/service"
	vs "github.com/postmanq/postmanq/module/validator/service"
	"go.uber.org/fx"
)

type PqModuleIn struct {
	fx.In
	ConfigProvider cs.ConfigProvider
	Validator      vs.Validator
}

type PqModuleOut struct {
	fx.Out
	Receiver interface{} `group:"component"`
}

func PqModule(params PqModuleIn) PqModuleOut {
	return PqModuleOut{
		Receiver: component.NewReceiver(
			params.ConfigProvider,
			qs.NewPool(),
			params.Validator,
		),
	}
}
