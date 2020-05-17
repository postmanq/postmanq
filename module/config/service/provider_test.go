package service_test

import (
	"github.com/postmanq/postmanq/module/config/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/config"
	"strings"
	"testing"
)

type A struct {
	B string
	C struct {
		D bool
		F int
	}
}

func TestConfigProviderSuite(t *testing.T) {
	suite.Run(t, new(ConfigProviderSuite))
}

type ConfigProviderSuite struct {
	suite.Suite
}

func (s *ConfigProviderSuite) TestPopulate() {
	reader := strings.NewReader("a: {b: bar, c: {d: true, f: 12}}")

	provider, err := service.NewConfigProviderByOptions(config.Source(reader))
	s.Nil(err)

	var a A
	s.Nil(provider.Populate("a", &a))
	s.Equal("bar", a.B)
	s.Equal(true, a.C.D)
	s.Equal(12, a.C.F)

	var f int
	s.Nil(provider.Populate("a.c.f", &f))
	s.Equal(12, f)
}
