package service_test

import (
	"github.com/postmanq/postmanq/mock/module/pipe/service/factory"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPipelineSuite(t *testing.T) {
	suite.Run(t, new(PipelineSuite))
}

type PipelineSuite struct {
	suite.Suite
	factory factory.PipelineFactory
}

func (s *PipelineSuite) SetupTest()           {}
func (s *PipelineSuite) TestFailureFactory()  {}
func (s *PipelineSuite) TestSuccessFactory()  {}
func (s *PipelineSuite) TestFailurePipeline() {}
func (s *PipelineSuite) TestSuccessPipeline() {}
