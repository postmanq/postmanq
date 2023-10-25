package temporal

import (
	"github.com/pkg/errors"
	"go.temporal.io/sdk/client"
	sdkworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"time"
)

var (
	ErrWorkflowNotFound = errors.New("api.workflow_not_found")
)

type Config struct {
	Address string `yaml:"address"`
}

type Context workflow.Context

type WorkerSettings = sdkworker.Options

type WorkerDescriptor struct {
	Workflow   WorkflowDescriptor
	Activities []ActivityDescriptor
}

type Signal struct {
	Channel string
	Message string
}

type WorkerFactoryIn struct {
	fx.In
	Client      Client
	Descriptors []WorkerDescriptor `group:"temporal_worker"`
}

type WorkflowType string

type ChildWorkflowOptions = workflow.ChildWorkflowOptions

type StartWorkflowOptions struct {
	client.StartWorkflowOptions
}

type WorkflowSettings struct {
	StartWorkflowOptions
	ChildWorkflowOptions
	ActivityOptions

	Delay       time.Duration
	ChunkSize   uint32
	StartSignal Signal
}

func (s *WorkflowSettings) GetStartWorkflowOptions() StartWorkflowOptions {
	return s.StartWorkflowOptions
}

type ActivityOptions workflow.ActivityOptions
