package main

import (
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/internal/services"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"smtp",
		fx.Provide(
			services.NewFxPluginDescriptor,
		),
	)
)
