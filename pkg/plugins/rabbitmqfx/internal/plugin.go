package internal

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"github.com/postmanq/postmanq/pkg/plugins/rabbitmqfx/rabbitmq"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func NewFxPluginDescriptor(
	executorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event],
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
					cfg:             cfg,
					conn:            conn,
					executorFactory: executorFactory,
				}, nil
			},
		},
	}
}

type plugin struct {
	cfg             rabbitmq.Config
	conn            *amqp.Connection
	executorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
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
			var event postmanqv1.Event
			err = proto.Unmarshal(delivery.Body, &event)
			if err != nil {
				return err
			}

			executor := p.executorFactory.Create(
				temporal.WithWorkflowType(temporal.WorkflowTypeSendEvent),
				temporal.WithWorkflowID(temporal.WorkflowTypeSendEvent, event.Uuid),
				temporal.WithWorkflowExecutionTimeout(time.Minute),
			)
			_, err = executor.Execute(ctx, &event)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			break
		}
	}
}
