package application

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type Application interface {
	Run(invoker interface{})
}

type PluginConstruct func() (Plugin, error)

type Plugin interface {
	GetKind() PluginKind
	OnReceive(ctx context.Context, next chan<- rxgo.Item) (rxgo.Observable, error)
	OnSend(ctx context.Context, item rxgo.Item) error
}
