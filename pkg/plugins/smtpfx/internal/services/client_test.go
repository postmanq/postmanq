package services_test

import (
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log_mock"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/internal/services"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

type ClientTestSuite struct {
	testutils.Suite
	builder smtp.ClientBuilder
}

func (s *ClientTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
	logger := log_mock.NewMockLogger(s.Ctrl)
	factory := services.NewFxClientBuilderFactory(logger)
	ip := "0.0.0.0"
	builder, err := factory.Create(s.Ctx, smtp.Config{
		IPs: []*string{&ip},
		Timeout: &smtp.TimeoutConfig{
			Connection: time.Minute,
			Hello:      time.Minute,
			Mail:       time.Minute,
			Rcpt:       time.Minute,
			Data:       time.Minute,
		},
	})
	s.Nil(err)
	s.builder = builder
}

func (s *ClientTestSuite) TestCreate() {
	_, err := s.builder.Create(s.Ctx, "smtp.google.com")
	s.Nil(err)
}

func (s *ClientTestSuite) TestSendEmail() {
	client, err := s.builder.Create(s.Ctx, "smtp.google.com")
	s.NotNil(client)
	s.Nil(err)

	err = client.Hello(s.Ctx, "mx.postmanq.io")
	s.Nil(err)

	err = client.Mail(s.Ctx, "test@postmanq.io")
	s.Nil(err)

	err = client.Rcpt(s.Ctx, "asolomonoff@gmail.com")
	s.Nil(err)
}
