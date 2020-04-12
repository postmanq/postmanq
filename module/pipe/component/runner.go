package component

import (
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/pipe/model"
)

type Runner struct {
	configProvider cs.ConfigProvider
	components     []interface{}
}

func NewRunner(
	configProvider cs.ConfigProvider,
	components []interface{},
) *Runner {
	return &Runner{
		configProvider: configProvider,
		components:     components,
	}
}

func (c *Runner) OnBootstrap() error {
	var cfg model.Config
	err := c.configProvider.Populate("pipes", &cfg)
	if err != nil {
		return err
	}

	return nil
}
