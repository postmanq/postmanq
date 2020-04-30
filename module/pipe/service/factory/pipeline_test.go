package factory_test

import (
	"github.com/postmanq/postmanq/mock/module/pipe/service/factory"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPipelineFactorySuite(t *testing.T) {
	suite.Run(t, new(PipelineFactorySuite))
}

type PipelineFactorySuite struct {
	suite.Suite
	factory factory.PipelineFactory
}
