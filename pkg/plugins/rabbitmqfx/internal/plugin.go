package internal

import (
	"context"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/plugins/rabbitmqfx/rabbitmq"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/reactivex/rxgo/v2"
)

func NewFxPluginDescriptor() postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name: "rabbitmq",
			Kind: postmanq.PluginKindReceiver,
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
					cfg:  cfg,
					conn: conn,
				}, nil
			},
		},
	}
}

type plugin struct {
	cfg  rabbitmq.Config
	conn *amqp.Connection
}

func (p plugin) OnReceive(ctx context.Context, next chan<- rxgo.Item) error {
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
			next <- rxgo.Of(delivery)
		case <-ctx.Done():
			break
		}
	}
}

func (p plugin) OnSend(ctx context.Context, item rxgo.Item) error {
	panic("implement me")
}
