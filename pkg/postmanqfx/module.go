package postmanqfx

import (
	"github.com/postmanq/postmanq/pkg/postmanqfx/internal/services"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"postmanq",
		fx.Provide(
			services.NewFxEventSenderFactory,
			services.NewFxInvoker,
		),
	)
)
