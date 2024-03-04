package postmanqfx

import (
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx"
	"github.com/postmanq/postmanq/pkg/postmanqfx/internal/services"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"postmanq",
		fx.Provide(
			services.NewFxEventSenderFactory,
			services.NewFxInvoker,
			temporalfx.NewFxWorkflowExecutorFactory[*postmanq.Event, *postmanq.Event],
			temporalfx.NewFxActivityExecutorFactory[*postmanq.Event, *postmanq.Event],
		),
	)
)
