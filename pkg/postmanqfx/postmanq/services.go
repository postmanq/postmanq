package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/reactivex/rxgo/v2"
)

type Postmanq interface {
}

type PluginConstruct func(provider config.Provider) (Plugin, error)

type Plugin interface {
	OnReceive(ctx context.Context, next chan<- rxgo.Item) error
	OnSend(ctx context.Context, item rxgo.Item) error
}
