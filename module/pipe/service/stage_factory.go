package service

import (
	"github.com/postmanq/postmanq/module/pipe/model"
)

type StageFactory interface {
	Create(model.Stage) Stage
}

func NewStageFactory() StageFactory {
	return &stageFactory{}
}

type stageFactory struct {
}

func (f *stageFactory) Create(cfg model.Stage) Stage {
	return nil
}
