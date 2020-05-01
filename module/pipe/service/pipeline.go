package service

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/service/factory"
	"github.com/postmanq/postmanq/module/pipe/service/stage"
)

type Pipeline interface {
	Init() error
	Run()
}

type pipeline struct {
	stages   []stage.Stage
	replicas int
}

func (p *pipeline) Init() error {
	for _, s := range p.stages {
		err := s.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *pipeline) Run() {
	var outs = make([]<-chan module.Delivery, len(p.stages))
	for i, s := range p.stages {
		for y := 0; y < p.replicas; y++ {
			var out <-chan module.Delivery

			if i == 0 {
				out = s.Start(nil)
			} else {
				out = s.Start(outs[i-1])
			}

			if y == 0 {
				outs[i] = out
			}
		}
	}
}

type PipelineFactory interface {
	Create(entity.Pipeline) (Pipeline, error)
}

func NewPipelineFactory(
	stageFactory factory.StageFactory,
) PipelineFactory {
	return &pipelineFactory{
		stageFactory: stageFactory,
	}
}

type pipelineFactory struct {
	stageFactory factory.StageFactory
}

func (f *pipelineFactory) Create(e entity.Pipeline) (Pipeline, error) {
	stages := make([]stage.Stage, len(e.Stages))
	for i, stageCfg := range e.Stages {
		s, err := f.stageFactory.Create(stageCfg)
		if err != nil {
			return nil, err
		}

		stages[i] = s
	}

	if e.Replicas == 0 {
		e.Replicas = 1
	}

	return &pipeline{
		stages:   stages,
		replicas: e.Replicas,
	}, nil
}
