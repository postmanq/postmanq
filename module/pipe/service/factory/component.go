package factory

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/config/service"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
)

type ComponentFactory interface {
	Register(module.ComponentDescriptor) error
	Create(*entity.Component) (interface{}, error)
}

func NewComponent(factory service.ConfigProviderFactory) ComponentFactory {
	return &componentFactory{
		descriptors: make(map[string]module.ComponentDescriptor),
		factory:     factory,
	}
}

type componentFactory struct {
	descriptors map[string]module.ComponentDescriptor
	factory     service.ConfigProviderFactory
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

	configProvider, err := f.factory.CreateFromMap(e.Config)
	if err != nil {
		return nil, err
	}

	return descriptor.Construct(configProvider), nil
}
