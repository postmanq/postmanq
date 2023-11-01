package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
	"go.temporal.io/sdk/workflow"
)

type Postmanq interface {
}

type PluginConstruct func(provider config.Provider) (Plugin, error)

type Plugin interface{}

type ReceiverPlugin interface {
	Plugin
	Receive(ctx context.Context) error
}

type WorkflowPlugin interface {
	Plugin
	GetActivityDescriptor() temporal.ActivityDescriptor
}

type EventSenderFactory interface {
	Create(pipeline Pipeline) EventSender
}

type EventSender interface {
	SendEvent(ctx workflow.Context, event *postmanqv1.Event) error
}
