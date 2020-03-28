package component

import (
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/pipe/model"
)

type Pipe struct {
	configProvider cs.ConfigProvider
	components     []interface{}
}

func NewPipe(
	configProvider cs.ConfigProvider,
	components []interface{},
) *Pipe {
	return &Pipe{
		configProvider: configProvider,
		components:     components,
	}
}

func (c *Pipe) OnBootstrap() error {
	var cfg model.Config
	err := c.configProvider.Populate("pipes", &cfg)
	if err != nil {
		return err
	}

	return nil
}
