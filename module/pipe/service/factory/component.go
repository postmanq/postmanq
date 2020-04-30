package factory

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
)

type ComponentFactory interface {
	Register(module.ComponentDescriptor) error
	Create(*entity.Component) (interface{}, error)
}

func NewComponentFactory() ComponentFactory {
	return &componentFactory{
		descriptors: make(map[string]module.ComponentDescriptor),
	}
}

type componentFactory struct {
	descriptors map[string]module.ComponentDescriptor
}

func (f *componentFactory) Register(descriptor module.ComponentDescriptor) error {
	_, ok := f.descriptors[descriptor.Name]
	if ok {
		return errors.ComponentDescriptorAlreadyDefined(descriptor.Name)
	}

	f.descriptors[descriptor.Name] = descriptor
	return nil
}

func (f *componentFactory) Create(e *entity.Component) (interface{}, error) {
	descriptor, ok := f.descriptors[e.Name]
	if !ok {
		return nil, errors.ComponentDescriptorNotDefined(e.Name)
	}

	return descriptor.Construct(e.Config), nil
}
