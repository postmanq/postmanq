package service

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/model"
)

type ComponentFactory interface {
	Register(module.ComponentDescriptor) error
	Create(model.Component) interface{}
}

func NewComponentFactory() ComponentFactory {
	return &componentFactory{}
}

type componentFactory struct {
}

func (f *componentFactory) Register(module.ComponentDescriptor) error {
	return nil
}

func (f *componentFactory) Create(model.Component) interface{} {
	return nil
}
