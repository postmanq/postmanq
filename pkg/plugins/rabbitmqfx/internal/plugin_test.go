package internal_test

import (
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPluginSuite(t *testing.T) {
	suite.Run(t, new(PluginSuite))
}

type PluginSuite struct {
	testutils.TemporalSuite
}

func (s *PluginSuite) Test() {
	s.True(s.ExecuteWorkflow())
	s.T().Log("finished test")
}
