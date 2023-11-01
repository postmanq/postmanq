package temporalfx

import (
	"github.com/postmanq/postmanq/pkg/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/temporalfx/internal"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"temporal",
		fx.Provide(
			internal.NewFxClient,
			internal.NewFxWorkerFactory,
			internal.NewFxWorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event],
			internal.NewFxActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event],
		),
	)
)
