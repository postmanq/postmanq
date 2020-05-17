package component

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/rabbitmq/entity"
	qs "github.com/postmanq/postmanq/module/rabbitmq/service"
	vs "github.com/postmanq/postmanq/module/validator/service"
	"github.com/streadway/amqp"
)

type receiver struct {
	ctx              context.Context
	configProvider   module.ConfigProvider
	pool             qs.Pool
	validator        vs.Validator
	repeatPublishers []qs.Publisher
	subscriber       qs.Subscriber
	url              string
	exchangeKind     string
}

func NewReceiver(
	ctx context.Context,
	pool qs.Pool,
	validator vs.Validator,
) module.ComponentOut {
	return module.ComponentOut{
		Descriptor: module.ComponentDescriptor{
			Name: "rabbitmq/receiver",
			Construct: func(configProvider module.ConfigProvider) interface{} {
				return &receiver{
					ctx:              ctx,
					configProvider:   configProvider,
					pool:             pool,
					validator:        validator,
					repeatPublishers: make([]qs.Publisher, 0),
					exchangeKind:     "direct",
				}
			},
		},
	}
}

func (c *receiver) OnInit() error {
	var cfg entity.Config
	err := c.configProvider.Populate("", &cfg)
	if err != nil {
		return err
	}

	err = c.validator.Struct(cfg)
	if err != nil {
		return err
	}

	if len(cfg.Prefix) == 0 {
		cfg.Prefix = "postmanq"
	}

	err = c.pool.Connect([]string{cfg.Url})
	if err != nil {
		return err
	}

	c.subscriber, err = c.pool.CreateSubscriber(c.ctx, entity.Queue{
		Name:          cfg.Prefix,
		Exchange:      cfg.Prefix,
		Durable:       true,
		PrefetchCount: 1,
	})
	if err != nil {
		return err
	}

	for _, repeat := range cfg.Repeats {
		repeatPublisher, err := c.pool.CreatePublisher(c.ctx, entity.Exchange{
			Name:    fmt.Sprintf("%s.repeat.%s", cfg.Prefix, repeat.String()),
			Kind:    c.exchangeKind,
			Durable: true,
			Args: amqp.Table{
				"x-message-ttl":          int64(repeat.Seconds()) * 1000,
				"x-dead-letter-exchange": cfg.Prefix,
			},
		})
		if err != nil {
			return err
		}

		c.repeatPublishers = append(c.repeatPublishers, repeatPublisher)
	}

	return nil
}

func (c *receiver) OnReceive(out chan module.Delivery) error {
	deliveries, err := c.subscriber.Subscribe(context.Background())
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		d := module.Delivery{}
		out <- d
		err = <-d.Err

		if err == nil {
			err := delivery.Ack(false)
			if err != nil {
				return err
			}
		} else {
			err := delivery.Reject(true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
