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
	Create(ctx context.Context, workflowType WorkflowType) (Worker, error)
	CreateByDescriptor(ctx context.Context, workerDescriptor WorkerDescriptor) (Worker, error)
}

type Worker interface {
	sdkworker.Worker
	RegisterWorkflowWithType(workflowType WorkflowType, i interface{})
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

type WorkflowDescriptor interface {
	GetWorkflowType() WorkflowType
}

type WorkflowOption interface {
	Apply(settings *WorkflowSettings)
}

func NewStartWorkflowOptions(options ...WorkflowOption) StartWorkflowOptions {
	settings := &WorkflowSettings{}

	for _, o := range options {
		o.Apply(settings)
	}

	return settings.GetStartWorkflowOptions()
}

type withWorkflowType string

func (w withWorkflowType) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.TaskQueue = string(w)
	o.ActivityOptions.TaskQueue = string(w)
}

func WithWorkflowType(workflowType WorkflowType) WorkflowOption {
	return withWorkflowType(workflowType)
}

// set child queue for split merge workflow

type withChildQueue string

func (w withChildQueue) Apply(o *WorkflowSettings) {
	o.ChildWorkflowOptions.TaskQueue = string(w)
}

func WithChildQueue(queue string) WorkflowOption {
	return withChildQueue(queue)
}

// set child queue for split merge workflow

type withWorkflowID string

func (w withWorkflowID) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.ID = string(w)
}

func WithWorkflowID(workflowType WorkflowType, args ...interface{}) WorkflowOption {
	if len(args) == 0 {
		return withWorkflowID(workflowType)
	}

	workflowIdParts := make([]string, len(args)+1)
	workflowIdParts[0] = string(workflowType)
	for i, arg := range args {
		workflowIdParts[i+1] = fmt.Sprint(arg)
	}

	return withWorkflowID(strings.Join(workflowIdParts, "_"))
}

// set schedule for cron job activities

type withCronSchedule string

func (w withCronSchedule) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.CronSchedule = string(w)
}

func WithCronSchedule(queue string) WorkflowOption {
	return withCronSchedule(queue)
}

type withWorkflowExecutionTimeout time.Duration

func (w withWorkflowExecutionTimeout) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.WorkflowExecutionTimeout = time.Duration(w)
}

func WithWorkflowExecutionTimeout(duration time.Duration) WorkflowOption {
	return withWorkflowExecutionTimeout(duration)
}

type withWorkflowTaskTimeout time.Duration

func (w withWorkflowTaskTimeout) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.WorkflowTaskTimeout = time.Duration(w)
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

type withActivityID string

func (w withActivityID) Apply(o *WorkflowSettings) {
	o.ActivityOptions.ActivityID = string(w)
}

func WithActivityID(activityID string) WorkflowOption {
	return withActivityID(activityID)
}

type withRetryPolicy temporal.RetryPolicy

func (w withRetryPolicy) Apply(o *WorkflowSettings) {
	policy := temporal.RetryPolicy(w)
	o.StartWorkflowOptions.RetryPolicy = &policy
}

func WithRetryPolicy(initialInterval time.Duration, maximumAttempts int32) WorkflowOption {
	return withRetryPolicy(temporal.RetryPolicy{
		InitialInterval: initialInterval,
		MaximumAttempts: maximumAttempts,
	})
}

type withChunkSize uint32

func (w withChunkSize) Apply(o *WorkflowSettings) {
	o.ChunkSize = uint32(w)
}

func WithChunkSize(chunkSize uint32) WorkflowOption {
	return withChunkSize(chunkSize)
}

type withStartSignal Signal

func (w withStartSignal) Apply(o *WorkflowSettings) {
	o.StartSignal = Signal(w)
}

type withStartDelay time.Duration

func (w withStartDelay) Apply(o *WorkflowSettings) {
	o.StartWorkflowOptions.StartDelay = time.Duration(w)
}

func WithStartDelay(duration time.Duration) WorkflowOption {
	return withStartDelay(duration)
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

type ActivityExecutorFactory[I any, O any] interface {
	Create(activityType string) ActivityExecutor[I, O]
}

type ActivityExecutor[I any, O any] interface {
	Execute(ctx workflow.Context, in I) (O, error)
}
