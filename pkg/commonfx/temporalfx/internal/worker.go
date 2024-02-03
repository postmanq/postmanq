package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"go.temporal.io/sdk/activity"
	sdkworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type worker struct {
	sdkworker.Worker
}

func (w *worker) RegisterWorkflowWithType(workflowType string, i interface{}) {
	w.RegisterWorkflowWithOptions(i, workflow.RegisterOptions{
		Name: string(workflowType),
	})
}

func (w *worker) RegisterActivityWithType(activityType string, i interface{}) {
	w.RegisterActivityWithOptions(i, activity.RegisterOptions{
		Name: activityType,
	})
}

func (w *worker) Run(interruptCh <-chan interface{}) error {
	err := w.Worker.Run(interruptCh)
	return err
}

func NewFxWorkerFactory(in temporal.WorkerFactoryIn) temporal.WorkerFactory {
	factory := &workerFactory{
		client:      in.Client,
		descriptors: collection.NewMap[string, temporal.WorkerDescriptor](),
	}
	for _, descriptor := range in.Descriptors {
		factory.descriptors.Set(descriptor.Workflow.Type, descriptor)
	}

	return factory
}

type workerFactory struct {
	client      temporal.Client
	descriptors collection.Map[string, temporal.WorkerDescriptor]
}

func (w workerFactory) Create(ctx context.Context, workflowType string) (temporal.Worker, error) {
	descriptor, ok := w.descriptors.Get(workflowType)
	if !ok {
		return nil, temporal.ErrWorkflowNotFound
	}

	return w.CreateByDescriptor(ctx, descriptor)
}

func (w workerFactory) CreateByDescriptor(ctx context.Context, descriptor temporal.WorkerDescriptor) (temporal.Worker, error) {
	wrk := &worker{
		Worker: sdkworker.New(w.client, descriptor.Workflow.Type, sdkworker.Options{
			Identity: fmt.Sprintf("%s.%s", descriptor.Workflow.Type, uuid.NewString()),
		}),
	}
	wrk.RegisterWorkflowWithType(descriptor.Workflow.Type, descriptor.Workflow.Func)
	for _, activityDescriptor := range descriptor.Activities {
		wrk.RegisterActivityWithType(activityDescriptor.Type, activityDescriptor.Func)
	}

	return wrk, nil
}
