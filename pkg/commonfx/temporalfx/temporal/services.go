package temporal

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	sdkworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"strings"
	"time"
)

type Client client.Client

type WorkerFactory interface {
	Create(ctx context.Context, workflowType string) (Worker, error)
	CreateByDescriptor(ctx context.Context, workerDescriptor WorkerDescriptor) (Worker, error)
}

type Worker interface {
	sdkworker.Worker
	RegisterWorkflowWithType(workflowType string, i interface{})
	RegisterActivityWithType(activityType string, i interface{})
}

type WorkerOption interface {
	Apply(*WorkerSettings)
}

type withDeadlockDetectionTimeout time.Duration

func (w withDeadlockDetectionTimeout) Apply(o *WorkerSettings) {
	o.DeadlockDetectionTimeout = time.Duration(w)
}

func WithDeadlockDetectionTimeout(duration time.Duration) WorkerOption {
	return withDeadlockDetectionTimeout(duration)
}

type workerOptionFunc func(options *WorkerSettings)

func (o workerOptionFunc) Apply(options *WorkerSettings) {
	o(options)
}

func WithEnableLoggingInReplay(val bool) WorkerOption {
	return workerOptionFunc(func(options *WorkerSettings) {
		options.EnableLoggingInReplay = val
	})
}

type WorkflowOption interface {
	Apply(settings *WorkflowSettings)
}

func NewStartWorkflowOptions(options ...WorkflowOption) StartWorkflowOptions {
	settings := &WorkflowSettings{}

	for _, o := range options {
		o.Apply(settings)
	}

	return settings.StartWorkflowOptions
}

func NewChildWorkflowOptions(options ...WorkflowOption) ChildWorkflowOptions {
	settings := &WorkflowSettings{}

	for _, o := range options {
		o.Apply(settings)
	}

	return settings.ChildWorkflowOptions
}

type withWorkflowType string

func (w withWorkflowType) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.TaskQueue = string(w)
	o.ChildWorkflowOptions.TaskQueue = string(w)
	o.ActivityOptions.TaskQueue = string(w)
}

func WithWorkflowType(workflowType string) WorkflowOption {
	return withWorkflowType(workflowType)
}

type withWorkflowId string

func (w withWorkflowId) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.ID = string(w)
	o.ChildWorkflowOptions.WorkflowID = string(w)
}

func WithWorkflowId(workflowType string, args ...interface{}) WorkflowOption {
	if len(args) == 0 {
		return withWorkflowId(workflowType)
	}

	workflowIdParts := make([]string, len(args)+1)
	workflowIdParts[0] = string(workflowType)
	for i, arg := range args {
		workflowIdParts[i+1] = fmt.Sprint(arg)
	}

	return withWorkflowId(strings.Join(workflowIdParts, "_"))
}

type withWorkflowExecutionTimeout time.Duration

func (w withWorkflowExecutionTimeout) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.WorkflowExecutionTimeout = time.Duration(w)
	o.ChildWorkflowOptions.WorkflowExecutionTimeout = time.Duration(w)
}

func WithWorkflowExecutionTimeout(duration time.Duration) WorkflowOption {
	return withWorkflowExecutionTimeout(duration)
}

type withWorkflowTaskTimeout time.Duration

func (w withWorkflowTaskTimeout) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.WorkflowTaskTimeout = time.Duration(w)
	o.ChildWorkflowOptions.WorkflowTaskTimeout = time.Duration(w)
}

func WithWorkflowTaskTimeout(duration time.Duration) WorkflowOption {
	return withWorkflowTaskTimeout(duration)
}

type withStartToCloseTimeout time.Duration

func (w withStartToCloseTimeout) Apply(o *WorkflowSettings) {
	o.ActivityOptions.StartToCloseTimeout = time.Duration(w)
}

func WithStartToCloseTimeout(duration time.Duration) WorkflowOption {
	return withStartToCloseTimeout(duration)
}

type withActivityId string

func (w withActivityId) Apply(o *WorkflowSettings) {
	o.ActivityOptions.ActivityID = string(w)
}

func WithActivityID(activityID string) WorkflowOption {
	return withActivityId(activityID)
}

type withRetryPolicy temporal.RetryPolicy

func (w withRetryPolicy) Apply(o *WorkflowSettings) {
	policy := temporal.RetryPolicy(w)
	o.StartWorkflowOptions.RetryPolicy = &policy
	o.ChildWorkflowOptions.RetryPolicy = &policy
	o.ActivityOptions.RetryPolicy = &policy
}

func WithRetryPolicy(initialInterval time.Duration, maximumAttempts int32) WorkflowOption {
	return withRetryPolicy(temporal.RetryPolicy{
		InitialInterval: initialInterval,
		MaximumAttempts: maximumAttempts,
	})
}

type withStartDelay time.Duration

func (w withStartDelay) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.StartDelay = time.Duration(w)
}

func WithStartDelay(duration time.Duration) WorkflowOption {
	return withStartDelay(duration)
}

func NewActivityOptions(options ...WorkflowOption) ActivityOptions {
	settings := &WorkflowSettings{}

	for _, o := range options {
		o.Apply(settings)
	}

	return settings.ActivityOptions
}

func WithActivityOptions(ctx workflow.Context, opts ActivityOptions) workflow.Context {
	if ctx == nil {
		return nil
	}

	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions(opts))
}

var InterruptCh = func() <-chan interface{} {
	return sdkworker.InterruptCh()
}

type WorkflowExecutorFactory[I any, O any] interface {
	Create(options ...WorkflowOption) WorkflowExecutor[I, O]
}

type WorkflowExecutor[I any, O any] interface {
	Execute(ctx context.Context, in I) (O, error)
}

type EventExecutor[I any, O any] interface {
	Execute(ctx workflow.Context, in I) (O, error)
}

type ActivityExecutorFactory[I any, O any] interface {
	Create(activityType string) EventExecutor[I, O]
}

type ChildWorkflowExecutorFactory[I any, O any] interface {
	Create(options ...WorkflowOption) EventExecutor[I, O]
}
