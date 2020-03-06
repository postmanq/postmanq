package component

import (
	"errors"
	cs "github.com/postmanq/postmanq/module/config/service"
	qs "github.com/postmanq/postmanq/module/queue/service"
	"go.uber.org/fx"
)

type ReceiverIn struct {
	fx.In
	ConfigProvider cs.ConfigProvider
	Pool           qs.Pool
}

type Receiver struct {
	configProvider cs.ConfigProvider
	pool           qs.Pool
}

func NewReceiver(
	configProvider cs.ConfigProvider,
	pool qs.Pool,
) *Receiver {
	return &Receiver{
		configProvider: configProvider,
		pool:           pool,
	}
}

func (c *Receiver) Init() error {
	var urls []string
	err := c.configProvider.Populate("queue", &urls)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return errors.New("queue.urls is not defined")
	}

	err = c.pool.Connect(urls)
	if err != nil {
		return err
	}

	return nil
}
