package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.temporal.io/sdk/workflow"
)

func NewFxPluginDescriptor(
	logger log.Logger,
	factory smtp.ClientBuilderFactory,
) postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name:       "smtp",
			Kind:       postmanq.PluginKindSender,
			MinVersion: 1.0,
			Construct: func(ctx context.Context, provider config.Provider) (postmanq.Plugin, error) {
				var cfg smtp.Config
				err := provider.Populate(&cfg)
				if err != nil {
					return nil, err
				}

				builder, err := factory.Create(ctx, cfg)
				if err != nil {
					return nil, err
				}

				return &plugin{
					cfg:     cfg,
					logger:  logger.Named("smtp_plugin"),
					builder: builder,
				}, nil
			},
		},
	}
}

type plugin struct {
	cfg     smtp.Config
	logger  log.Logger
	builder smtp.ClientBuilder
}

func (p plugin) GetType() string {
	return "ActivityTypeSMTP"
}

func (p plugin) OnEvent(ctx workflow.Context, event *postmanqv1.Event) (*postmanqv1.Event, error) {
	//TODO implement me
	panic("implement me")
}
