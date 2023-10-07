package services_test

import (
	"github.com/postmanq/postmanq/pkg/configfx/internal/services"
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
	factory := services.NewFxConfigProviderFactory()
	provider, err := factory.Create(config.Static(map[string]string{
		"b": "hello",
	}))
	s.Nil(err)

	var a A
	s.Nil(provider.Populate(&a))
	s.Equal("hello", a.B)

	reader := strings.NewReader("a: {b: bar, c: {d: true, f: 12}}")
	provider, err = factory.Create(config.Source(reader))
	s.Nil(err)

	s.Nil(provider.PopulateByKey("a", &a))
	s.Equal("bar", a.B)
	s.Equal(true, a.C.D)
	s.Equal(12, a.C.F)

	var f int
	s.Nil(provider.PopulateByKey("a.c.f", &f))
	s.Equal(12, f)
}
