package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"go.temporal.io/sdk/workflow"
)

type PluginConstruct func(ctx context.Context, provider config.Provider) (Plugin, error)

type Plugin interface{}

type ReceiverPlugin interface {
	Plugin
	Receive(ctx context.Context) error
}

type WorkflowPlugin interface {
	Plugin
	GetType() string
	OnEvent(ctx workflow.Context, event *postmanqv1.Event) (*postmanqv1.Event, error)
}

type EventSenderFactory interface {
	Create(pipeline *Pipeline) EventSender
}

type EventSender interface {
	SendEvent(ctx workflow.Context, event *postmanqv1.Event) (*postmanqv1.Event, error)
}

type SendEventWorkflow func(ctx workflow.Context, event *postmanqv1.Event) (*postmanqv1.Event, error)

func (w SendEventWorkflow) GetWorkflowType() temporal.WorkflowType {
	return temporal.WorkflowTypeSendEvent
}

type Invoker interface {
	Configure(ctx context.Context) error
	Run(ctx context.Context) error
}
