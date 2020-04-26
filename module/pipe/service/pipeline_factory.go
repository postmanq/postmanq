package service

import "github.com/postmanq/postmanq/module/pipe/model"

type PipelineFactory interface {
	Create(model.Pipeline) (Pipeline, error)
}

func NewPipelineFactory(
	stageFactory StageFactory,
	componentFactory ComponentFactory,
) PipelineFactory {
	return &pipelineFactory{}
}

type pipelineFactory struct {
}

func (f *pipelineFactory) Create(cfg model.Pipeline) (Pipeline, error) {
	return nil, nil
}
