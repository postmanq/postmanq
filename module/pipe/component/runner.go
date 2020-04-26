package component

import (
	"github.com/postmanq/postmanq/module"
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/pipe/model"
	"github.com/postmanq/postmanq/module/pipe/service"
	"go.uber.org/fx"
)

type RunnerIn struct {
	fx.In
	ConfigProvider   cs.ConfigProvider
	ComponentFactory service.ComponentFactory
	PipelineFactory  service.PipelineFactory
	Components       []module.ComponentDescriptor `group:"component"`
}

type Runner struct {
	configProvider   cs.ConfigProvider
	components       []module.ComponentDescriptor
	componentFactory service.ComponentFactory
	pipelineFactory  service.PipelineFactory
}

func NewRunner(in RunnerIn) *Runner {
	return &Runner{
		configProvider: in.ConfigProvider,
		components:     in.Components,
	}
}

func (c *Runner) Run() error {
	var pipelineConfigs []model.Pipeline
	err := c.configProvider.Populate("pipelines", &pipelineConfigs)
	if err != nil {
		return err
	}

	for _, component := range c.components {
		err = c.componentFactory.Register(component)
		if err != nil {
			return err
		}
	}

	pipelines := make([]service.Pipeline, 0)
	for _, pipelineConfig := range pipelineConfigs {
		pipeline, err := c.pipelineFactory.Create(pipelineConfig)
		if err != nil {
			return err
		}

		pipelines = append(pipelines, pipeline)
	}

	for _, pipeline := range pipelines {
		go pipeline.Run()
	}

	return nil
}
