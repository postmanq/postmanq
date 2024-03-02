package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/plugins/serverfx/server"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
)

func NewFxPluginDescriptor(
	ctx context.Context,
	factory server.Factory,
	eventServiceServer postmanqv1.EventServiceServer,
) postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name:       "server",
			Kind:       postmanq.PluginKindReceiver,
			MinVersion: 1.0,
			Construct: func(ctx context.Context, pipeline postmanq.Pipeline, provider config.Provider) (postmanq.Plugin, error) {
				var cfg server.Config
				err := provider.Populate(&cfg)
				if err != nil {
					return nil, err
				}

				srv, err := factory.Create(ctx, cfg)
				if err != nil {
					return nil, err
				}

				return &plugin{
					server: srv,
					descriptor: server.Descriptor{
						Server:               eventServiceServer,
						GRPCGatewayRegistrar: postmanqv1.RegisterEventServiceHandlerFromEndpoint,
					},
				}, nil
			},
		},
	}
}

type plugin struct {
	server     server.Server
	descriptor server.Descriptor
}

func (p *plugin) Receive(ctx context.Context) error {
	err := p.server.Register(p.descriptor)
	if err != nil {
		return err
	}

	return p.server.Start()
}
