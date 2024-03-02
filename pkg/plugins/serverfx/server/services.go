package server

import (
	"context"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
)

type Factory interface {
	Create(ctx context.Context, cfg Config) (Server, error)
}

type Server interface {
	Register(descriptor Descriptor) error
	Start() error
}

type EventServiceServerFactory interface {
	Create(ctx context.Context, pipeline postmanq.Pipeline) EventServiceServer
}
