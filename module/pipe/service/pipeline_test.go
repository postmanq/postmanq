package service_test

import (
	"fmt"
	sf "github.com/postmanq/postmanq/mock/module/pipe/service/factory"
	"github.com/postmanq/postmanq/mock/module/pipe/service/stage"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"github.com/postmanq/postmanq/module/pipe/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPipelineSuite(t *testing.T) {
	suite.Run(t, new(PipelineSuite))
}

type PipelineSuite struct {
	suite.Suite
	pipelineFactory service.PipelineFactory
	stageFactory    *sf.StageFactory
}

func (s *PipelineSuite) SetupTest() {
	s.stageFactory = new(sf.StageFactory)
	s.pipelineFactory = service.NewPipelineFactory(
		s.stageFactory,
	)
}

func (s *PipelineSuite) TestFailureFactory() {
	st := &entity.Stage{
		Type: "wrong",
	}
	expectedErr := errors.ConstructNotDefinedForStage(st)
	s.stageFactory.On("Create", mock.Anything).Return(nil, expectedErr).Once()
	pipeline, err := s.pipelineFactory.Create(entity.Pipeline{
		Stages: []*entity.Stage{st},
	})
	s.Nil(pipeline)
	s.NotNil(err)
	s.Equal(expectedErr, err)
}

func (s *PipelineSuite) TestSuccessFactory() {
	s.stageFactory.On("Create", mock.Anything).Return(new(stage.Stage), nil).Once()
	pipeline, err := s.pipelineFactory.Create(entity.Pipeline{
		Stages: []*entity.Stage{{}},
	})
	s.NotNil(pipeline)
	s.Nil(err)
}

func (s *PipelineSuite) TestFailurePipeline() {
	expectedErr := fmt.Errorf("stage error")
	st := new(stage.Stage)
	st.On("Init").Return(expectedErr).Once()
	s.stageFactory.On("Create", mock.Anything).Return(st, nil).Once()

	pipeline, err := s.pipelineFactory.Create(entity.Pipeline{
		Stages: []*entity.Stage{{}},
	})
	s.NotNil(pipeline)
	s.Nil(err)

	err = pipeline.Init()
	s.NotNil(err)
	s.Equal(expectedErr, err)
}

func (s *PipelineSuite) TestSuccessPipeline() {
	st := new(stage.Stage)
	st.On("Init").Return(nil).Once()
	st.On("Start", mock.Anything).Return(nil)
	s.stageFactory.On("Create", mock.Anything).Return(st, nil).Once()

	pipeline, err := s.pipelineFactory.Create(entity.Pipeline{
		Replicas: 4,
		Stages:   []*entity.Stage{{}},
	})
	s.NotNil(pipeline)
	s.Nil(err)

	err = pipeline.Init()
	s.Nil(err)
	st.AssertNumberOfCalls(s.T(), "Init", 1)

	pipeline.Run()
	st.AssertNumberOfCalls(s.T(), "Start", 4)

}
