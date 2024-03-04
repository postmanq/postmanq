package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"go.temporal.io/sdk/workflow"
)

func NewWorkflowExecutorFactory[I any, O any](temporalClient temporal.Client) temporal.WorkflowExecutorFactory[I, O] {
	return &workflowExecutorFactory[I, O]{
		cl: temporalClient,
	}
}

type workflowExecutorFactory[I any, O any] struct {
	cl temporal.Client
}

func (e *workflowExecutorFactory[I, O]) Create(options ...temporal.WorkflowOption) temporal.WorkflowExecutor[I, O] {
	return &workflowExecutor[I, O]{
		cl:      e.cl,
		options: options,
	}
}

type workflowExecutor[I any, O any] struct {
	cl      temporal.Client
	options []temporal.WorkflowOption
}

func (e *workflowExecutor[I, O]) Execute(ctx context.Context, in I) (O, error) {
	var out O
	workflowOptions := temporal.NewStartWorkflowOptions(e.options...)

	wr, err := e.cl.ExecuteWorkflow(
		ctx,
		workflowOptions.StartWorkflowOptions,
		workflowOptions.StartWorkflowOptions.TaskQueue,
		in,
	)
	if err != nil {
		return out, err
	}

	err = wr.Get(ctx, &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

func NewActivityExecutorFactory[I any, O any]() temporal.ActivityExecutorFactory[I, O] {
	return &activityExecutorFactory[I, O]{}
}

type activityExecutorFactory[I any, O any] struct{}

func (a activityExecutorFactory[I, O]) Create(activityType string) temporal.ActivityExecutor[I, O] {
	return &activityExecutor[I, O]{
		activityType: activityType,
	}
}

type activityExecutor[I any, O any] struct {
	activityType string
}

func (a *activityExecutor[I, O]) Execute(ctx workflow.Context, in I) (O, error) {
	var out O
	err := workflow.ExecuteActivity(ctx, a.activityType, in).Get(ctx, &out)
	return out, err
}
