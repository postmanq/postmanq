package component

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/module"
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/queue/model"
	qs "github.com/postmanq/postmanq/module/queue/service"
	vs "github.com/postmanq/postmanq/module/validator/service"
	"github.com/streadway/amqp"
)

const (
	prefix       = "postmanq"
	exchangeKind = "direct"
)

type Receiver struct {
	configProvider   cs.ConfigProvider
	pool             qs.Pool
	validator        vs.Validator
	repeatPublishers []qs.Publisher
	subscriber       qs.Subscriber
}

func NewReceiver(
	configProvider cs.ConfigProvider,
	pool qs.Pool,
	validator vs.Validator,
) *Receiver {
	return &Receiver{
		configProvider:   configProvider,
		pool:             pool,
		validator:        validator,
		repeatPublishers: make([]qs.Publisher, 0),
	}
}

func (c *Receiver) OnInit() error {
	var cfg model.Config
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

	c.subscriber, err = c.pool.CreateSubscriber(context.Background(), model.Queue{
		Name:          prefix,
		Exchange:      prefix,
		Durable:       true,
		PrefetchCount: cfg.PrefetchMessageCount,
	})
	if err != nil {
		return err
	}

	for _, repeat := range cfg.Repeats {
		repeatPublisher, err := c.pool.CreatePublisher(context.Background(), model.Exchange{
			Name:    fmt.Sprintf("%s.repeat.%s", prefix, repeat.String()),
			Kind:    exchangeKind,
			Durable: true,
			Args: amqp.Table{
				"x-message-ttl":          int64(repeat.Seconds()) * 1000,
				"x-dead-letter-exchange": prefix,
			},
		})
		if err != nil {
			return err
		}

		c.repeatPublishers = append(c.repeatPublishers, repeatPublisher)
	}

	return nil
}

func (c *Receiver) OnReceive(deliveries chan module.Delivery, results chan module.Delivery) error {
	_, err := c.subscriber.Subscribe(context.Background())
	if err != nil {
		return err
	}

	return nil
}
