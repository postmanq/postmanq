package component

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/module"
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/rabbitmq/entity"
	qs "github.com/postmanq/postmanq/module/rabbitmq/service"
	vs "github.com/postmanq/postmanq/module/validator/service"
	"github.com/streadway/amqp"
)

type Receiver interface {
	OnInit() error
	OnReceive(out chan module.Delivery) error
}

type receiver struct {
	configProvider   cs.ConfigProvider
	pool             qs.Pool
	validator        vs.Validator
	repeatPublishers []qs.Publisher
	subscriber       qs.Subscriber
	prefix           string
	exchangeKind     string
}

func NewReceiver(
	configProvider cs.ConfigProvider,
	pool qs.Pool,
	validator vs.Validator,
) module.ComponentOut {
	return module.ComponentOut{
		Descriptor: module.ComponentDescriptor{
			Name: "rabbitmq/receiver",
			Construct: func(cfg module.ComponentConfig) interface{} {
				var prefix string
				if rawPrefix, ok := cfg["prefix"]; ok {
					prefix = rawPrefix.(string)
				}

				return &receiver{
					configProvider:   configProvider,
					pool:             pool,
					validator:        validator,
					repeatPublishers: make([]qs.Publisher, 0),
					prefix:           prefix,
					exchangeKind:     "direct",
				}
			},
		},
	}
}

func (c *receiver) OnInit() error {
	var cfg entity.Config
	err := c.configProvider.Populate("queue", &cfg)
	if err != nil {
		return err
	}

	err = c.validator.Struct(cfg)
	if err != nil {
		return err
	}

	err = c.pool.Connect([]string{cfg.Url})
	if err != nil {
		return err
	}

	c.subscriber, err = c.pool.CreateSubscriber(context.Background(), entity.Queue{
		Name:          c.prefix,
		Exchange:      c.prefix,
		Durable:       true,
		PrefetchCount: 1,
	})
	if err != nil {
		return err
	}

	for _, repeat := range cfg.Repeats {
		repeatPublisher, err := c.pool.CreatePublisher(context.Background(), entity.Exchange{
			Name:    fmt.Sprintf("%s.repeat.%s", c.prefix, repeat.String()),
			Kind:    c.exchangeKind,
			Durable: true,
			Args: amqp.Table{
				"x-message-ttl":          int64(repeat.Seconds()) * 1000,
				"x-dead-letter-exchange": c.prefix,
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
