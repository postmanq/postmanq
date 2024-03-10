package temporalfx

import (
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/internal/services"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"temporal",
		fx.Provide(
			services.NewFxClient,
			services.NewFxWorkerFactory,
		),
	)
)

func NewFxWorkflowExecutorFactory[I any, O any](temporalClient temporal.Client) temporal.WorkflowExecutorFactory[I, O] {
	return services.NewWorkflowExecutorFactory[I, O](temporalClient)
}

func NewFxChildWorkflowExecutorFactory[I any, O any]() temporal.ChildWorkflowExecutorFactory[I, O] {
	return services.NewChildWorkflowExecutorFactory[I, O]()
}

func NewFxActivityExecutorFactory[I any, O any]() temporal.ActivityExecutorFactory[I, O] {
	return services.NewActivityExecutorFactory[I, O]()
}
