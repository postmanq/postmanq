package temporalfx

import (
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/internal"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"temporal",
		fx.Provide(
			internal.NewFxClient,
			internal.NewFxWorkerFactory,
			internal.NewFxWorkflowExecutorFactory[*postmanq.Event, *postmanq.Event],
			internal.NewFxActivityExecutorFactory[*postmanq.Event, *postmanq.Event],
		),
	)
)
