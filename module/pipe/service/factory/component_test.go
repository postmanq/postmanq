package factory_test

import (
	mm "github.com/postmanq/postmanq/mock/module"
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"github.com/postmanq/postmanq/module/pipe/service/factory"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestComponentFactorySuite(t *testing.T) {
	suite.Run(t, new(ComponentFactorySuite))
}

type ComponentFactorySuite struct {
	suite.Suite
	factory         factory.ComponentFactory
	validDescriptor module.ComponentDescriptor
}

func (s *ComponentFactorySuite) SetupTest() {
	s.factory = factory.NewComponent()
	s.validDescriptor = module.ComponentDescriptor{
		Name: "component1",
		Construct: func(configs module.ComponentConfig) interface{} {
			return new(mm.InitComponent)
		},
	}
}

func (s *ComponentFactorySuite) TestFailure() {
	s.Nil(s.factory.Register(s.validDescriptor))
	err := s.factory.Register(s.validDescriptor)
	s.NotNil(err)
	s.Equal(errors.ComponentDescriptorAlreadyDefined(s.validDescriptor.Name), err)

	cfg := &entity.Component{
		Name: "component2",
	}
	comp, err := s.factory.Create(cfg)
	s.Nil(comp)
	s.NotNil(err)
	s.Equal(errors.ComponentDescriptorNotDefined(cfg.Name), err)
}

func (s *ComponentFactorySuite) TestSuccess() {
	s.Nil(s.factory.Register(s.validDescriptor))
	comp, err := s.factory.Create(&entity.Component{
		Name: "component1",
	})
	s.NotNil(comp)
	s.Nil(err)
	s.Implements(new(module.InitComponent), comp)
}
