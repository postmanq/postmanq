package server

import "context"

type Factory interface {
	Create(ctx context.Context, cfg Config) (Server, error)
}

type Server interface {
	Register(descriptor Descriptor) error
	Start() error
}
