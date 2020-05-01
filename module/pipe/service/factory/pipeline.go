package factory

import (
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/service"
)

type PipelineFactory interface {
	Create(entity.Pipeline) (service.Pipeline, error)
}

func NewPipelineFactory(
	stageFactory StageFactory,
	componentFactory ComponentFactory,
) PipelineFactory {
	return &pipelineFactory{}
}

type pipelineFactory struct {
}

func (f *pipelineFactory) Create(cfg entity.Pipeline) (service.Pipeline, error) {

	return nil, nil
}
