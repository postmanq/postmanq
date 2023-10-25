package rabbitmqfx

import (
	"github.com/postmanq/postmanq/pkg/plugins/rabbitmqfx/internal"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"rabbitmq",
		fx.Provide(
			internal.NewFxPluginDescriptor,
		),
	)
)
