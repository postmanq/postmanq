package component

import (
	"github.com/postmanq/postmanq/module"
	cs "github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/service"
	"github.com/postmanq/postmanq/module/pipe/service/factory"
	"go.uber.org/fx"
)

type RunnerIn struct {
	fx.In
	ConfigProvider   cs.ConfigProvider
	ComponentFactory factory.ComponentFactory
	PipelineFactory  service.PipelineFactory
	Descriptors      []module.ComponentDescriptor `group:"component"`
}

type Runner struct {
	configProvider   cs.ConfigProvider
	descriptors      []module.ComponentDescriptor
	componentFactory factory.ComponentFactory
	pipelineFactory  service.PipelineFactory
}

func NewRunner(in RunnerIn) *Runner {
	return &Runner{
		configProvider: in.ConfigProvider,
		descriptors:    in.Descriptors,
	}
}

func (c *Runner) Run() error {
	var pipelineConfigs []entity.Pipeline
	err := c.configProvider.Populate("pipelines", &pipelineConfigs)
	if err != nil {
		return err
	}

	for _, component := range c.descriptors {
		err = c.componentFactory.Register(component)
		if err != nil {
			return err
		}
	}

	pipelines := make([]service.Pipeline, len(pipelineConfigs))
	for i, pipelineConfig := range pipelineConfigs {
		pipeline, err := c.pipelineFactory.Create(pipelineConfig)
		if err != nil {
			return err
		}

		pipelines[i] = pipeline
	}

	for _, pipeline := range pipelines {
		err := pipeline.Init()
		if err != nil {
			return err
		}
	}

	for _, pipeline := range pipelines {
		go pipeline.Run()
	}

	return nil
}
