package main

import (
	"github.com/postmanq/postmanq/pkg/plugins/serverfx/internal/services"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"server",
		fx.Provide(
			services.NewFxEventServiceServerFactory,
			services.NewFxUnionServerFactory,
			services.NewFxPluginDescriptor,
		),
	)
)
