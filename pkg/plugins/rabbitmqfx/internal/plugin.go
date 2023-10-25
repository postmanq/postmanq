package internal

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/plugins/rabbitmqfx/rabbitmq"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewFxPluginDescriptor(
	executor temporal.WorkflowExecutor[*postmanqv1.Event, *postmanqv1.Event],
) postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name:       "rabbitmq",
			Kind:       postmanq.PluginKindReceiver,
			MinVersion: 1.0,
			Construct: func(provider config.Provider) (postmanq.Plugin, error) {
				var cfg rabbitmq.Config
				err := provider.Populate(&cfg)
				if err != nil {
					return nil, err
				}

				conn, err := amqp.Dial(cfg.Url)
				if err != nil {
					return nil, err
				}

				return &plugin{
					cfg:      cfg,
					conn:     conn,
					executor: executor,
				}, nil
			},
		},
	}
}

type plugin struct {
	cfg      rabbitmq.Config
	conn     *amqp.Connection
	executor temporal.WorkflowExecutor[*postmanqv1.Event, *postmanqv1.Event]
}

func (p plugin) Receive(ctx context.Context) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	deliveries, err := ch.ConsumeWithContext(
		ctx,
		p.cfg.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case delivery := <-deliveries:
			_, err := p.executor.Execute(ctx, nil)
			if err != nil {
				return err
			}

			var event postmanqv1.Event
			err = proto.Unmarshal(delivery.Body, &event)
			if err != nil {
				return err
			}

			_, err = p.executor.Execute(ctx, &event)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			break
		}
	}
}
