package postmanqfx

import (
	temporal "github.com/postmanq/postmanq/pkg/temporalfx/internal"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"postmanq",
		fx.Provide(
			temporal.NewFxClient,
			temporal.NewFxWorkerFactory,
		),
	)
)
