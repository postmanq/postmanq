package services_test

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/internal/services"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestEmailParserTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(EmailParserTestSuite))
}

type EmailParserTestSuite struct {
	testutils.Suite
	emailParser smtp.EmailParser
}

func (s *EmailParserTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.emailParser = services.NewFxEmailParser()
}

func (s *EmailParserTestSuite) TestParse() {
	email, err := s.emailParser.Parse(randomdata.Address())
	s.Nil(email)
	s.NotNil(err)

	email, err = s.emailParser.Parse(randomdata.Email())
	s.IsType(&smtp.EmailAddress{}, email)
	s.Nil(err)
}
