package main

import (
	"github.com/postmanq/postmanq/pkg/plugins/templatefx/internal/services"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"template",
		fx.Provide(
			services.NewFxPluginDescriptor,
		),
	)
)
