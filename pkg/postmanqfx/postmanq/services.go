package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"go.temporal.io/sdk/workflow"
)

type PluginConstruct func(ctx context.Context, pipeline Pipeline, provider config.Provider) (Plugin, error)

type Plugin interface{}

type ReceiverPlugin interface {
	Plugin
	Receive(ctx context.Context) error
}

type WorkflowPlugin interface {
	Plugin
	GetType() string
	OnEvent(ctx context.Context, event *Event) (*Event, error)
}

type EventSenderFactory interface {
	Create(pipeline *Pipeline) EventSender
}

type EventSender interface {
	SendEvent(ctx workflow.Context, event *Event) (*Event, error)
}

type Invoker interface {
	Configure(ctx context.Context) error
	Run(ctx context.Context) error
}
