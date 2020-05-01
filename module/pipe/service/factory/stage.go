package factory

import (
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"github.com/postmanq/postmanq/module/pipe/service/stage"
	"go.uber.org/fx"
)

type StageFactory interface {
	Create(*entity.Stage) (stage.Stage, error)
}

type StageFactoryIn struct {
	fx.In
	ComponentFactory ComponentFactory
	Stages           []stage.Descriptor `group:"stage"`
}

func NewStage(in StageFactoryIn) StageFactory {
	descriptors := make(map[string]stage.Descriptor)
	for _, descriptor := range in.Stages {
		descriptors[descriptor.Name] = descriptor
	}

	return &stageFactory{
		componentFactory: in.ComponentFactory,
		descriptors:      descriptors,
	}
}

type stageFactory struct {
	componentFactory ComponentFactory
	descriptors      map[string]stage.Descriptor
}

func (f *stageFactory) Create(e *entity.Stage) (stage.Stage, error) {
	descriptor, ok := f.descriptors[e.Type]
	if !ok {
		return nil, errors.UnknownStageType(e)
	}

	switch descriptor.Type {
	case stage.ArgTypeSingle:
		if e.Component == nil {
			return nil, errors.ComponentNotDefinedForStage(e)
		}

		component, err := f.componentFactory.Create(e.Component)
		if err != nil {
			return nil, err
		}

		return descriptor.Constructor(e, component)
	case stage.ArgTypeMulti:
		componentsLen := len(e.Components)
		if componentsLen == 0 {
			return nil, errors.ComponentsNotDefinedForStage(e)
		}

		components := make([]interface{}, componentsLen)
		for i, item := range e.Components {
			component, err := f.componentFactory.Create(item)
			if err != nil {
				return nil, err
			}

			components[i] = component
		}

		return descriptor.Constructor(e, components)
	}

	return nil, errors.ConstructNotDefinedForStage(e)
}
