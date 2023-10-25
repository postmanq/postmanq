package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/configfx/config"
)

type Postmanq interface {
}

type PluginConstruct func(provider config.Provider) (Plugin, error)

type Plugin interface{}

type ReceiverPlugin interface {
	Plugin
	Receive(ctx context.Context) error
}

type MiddlewarePlugin interface {
	Plugin
	Next(ctx context.Context, event *Event) error
}

type SenderPlugin interface {
	Plugin
	Send(ctx context.Context, event *Event) error
}
