package factory_test

import (
	mf "github.com/postmanq/postmanq/mock/module/pipe/service/factory"
	ms "github.com/postmanq/postmanq/mock/module/pipe/service/stage"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"github.com/postmanq/postmanq/module/pipe/service/factory"
	"github.com/postmanq/postmanq/module/pipe/service/stage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type stageDataItem struct {
	stage     *entity.Stage
	construct func(*entity.Stage) error
}

var (
	failureStageDataItems = []stageDataItem{
		{
			stage: &entity.Stage{
				Type: "wrong1",
			},
			construct: errors.UnknownStageType,
		},
		{
			stage: &entity.Stage{
				Type: "wrong2",
			},
			construct: errors.ConstructNotDefinedForStage,
		},
		{
			stage: &entity.Stage{
				Type: "single",
			},
			construct: errors.ComponentNotDefinedForStage,
		},
		{
			stage: &entity.Stage{
				Type: "multi",
			},
			construct: errors.ComponentsNotDefinedForStage,
		},
	}
	successStageDataItems = []stageDataItem{
		{
			stage: &entity.Stage{
				Name: "stage1",
				Type: "single",
				Component: &entity.Component{
					Name: "component1",
				},
			},
		},
		{
			stage: &entity.Stage{
				Name: "stage2",
				Type: "multi",
				Components: []*entity.Component{
					{
						Name: "component1",
					},
				},
			},
		},
	}
)

func TestStageFactorySuite(t *testing.T) {
	suite.Run(t, new(StageFactorySuite))
}

type StageFactorySuite struct {
	suite.Suite
	stageFactory     factory.StageFactory
	componentFactory *mf.ComponentFactory
}

func (s *StageFactorySuite) SetupTest() {
	s.componentFactory = new(mf.ComponentFactory)
	s.stageFactory = factory.NewStageFactory(
		factory.StageFactoryIn{
			ComponentFactory: s.componentFactory,
			Stages: []stage.Descriptor{
				{
					Name: "single",
					Type: stage.SingleComponentType,
					Constructor: func(e *entity.Stage, c interface{}) (stage.Stage, error) {
						return new(ms.Stage), nil
					},
				},
				{
					Name: "multi",
					Type: stage.MultiComponentType,
					Constructor: func(e *entity.Stage, c interface{}) (stage.Stage, error) {
						return new(ms.Stage), nil
					},
				},
				{
					Name: "wrong2",
				},
			},
		},
	)
}

func (s *StageFactorySuite) TestFailure() {
	for _, failureDataItem := range failureStageDataItems {
		s.componentFactory.On("Create", mock.Anything).Return(nil, nil).Once()
		st, err := s.stageFactory.Create(failureDataItem.stage)
		s.Nil(st)
		s.NotNil(err)
		s.Equal(failureDataItem.construct(failureDataItem.stage), err)
	}
}

func (s *StageFactorySuite) TestSuccess() {
	for _, item := range successStageDataItems {
		s.componentFactory.On("Create", mock.Anything).Return(new(ms.Stage), nil).Once()
		st, err := s.stageFactory.Create(item.stage)
		s.Nil(err)
		s.NotNil(st)
		s.Implements(new(stage.Stage), st)
	}
}
