package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/postmanq/postmanq/pkg/collection"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
	"go.temporal.io/sdk/activity"
	sdkworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type worker struct {
	sdkworker.Worker
}

func (w *worker) RegisterWorkflowWithType(workflowType temporal.WorkflowType, i interface{}) {
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
		descriptors: collection.NewMap[temporal.WorkflowType, temporal.WorkerDescriptor](),
	}
	for _, descriptor := range in.Descriptors {
		factory.descriptors.Set(descriptor.Workflow.GetWorkflowType(), descriptor)
	}

	return factory
}

type workerFactory struct {
	client      temporal.Client
	descriptors collection.Map[temporal.WorkflowType, temporal.WorkerDescriptor]
}

func (w workerFactory) Create(ctx context.Context, workflowType temporal.WorkflowType) (temporal.Worker, error) {
	descriptor, ok := w.descriptors.Get(workflowType)
	if !ok {
		return nil, temporal.ErrWorkflowNotFound
	}

	return w.CreateByDescriptor(ctx, descriptor)
}

func (w workerFactory) CreateByDescriptor(ctx context.Context, workerDescriptor temporal.WorkerDescriptor) (temporal.Worker, error) {
	workflowType := workerDescriptor.Workflow.GetWorkflowType()
	wrk := &worker{
		Worker: sdkworker.New(w.client, string(workflowType), sdkworker.Options{
			Identity: fmt.Sprintf("%s.%s", workflowType, uuid.NewString()),
		}),
	}
	wrk.RegisterWorkflowWithType(workflowType, workerDescriptor.Workflow)
	for _, activityDescriptor := range workerDescriptor.Activities {
		wrk.RegisterActivityWithType(activityDescriptor.GetActivityType(), activityDescriptor)
	}

	return wrk, nil
}
