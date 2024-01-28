package services_test

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log_mock"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/internal/services"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestMxResolverTestSuite(t *testing.T) {
	suite.Run(t, new(MxResolverTestSuite))
}

type MxResolverTestSuite struct {
	testutils.Suite
	mxResolver smtp.MxResolver
}

func (s *MxResolverTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
	logger := log_mock.NewMockLogger(s.Ctrl)
	logger.EXPECT().Named(gomock.Any()).AnyTimes().Return(logger)
	logger.EXPECT().WithCtx(gomock.Any()).AnyTimes().Return(logger)
	logger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes().Return()
	logger.EXPECT().Debugf(gomock.Any()).AnyTimes().Return()
	s.mxResolver = services.NewFxMxResolver(logger)
}

func (s *MxResolverTestSuite) TestResolve() {
	sl, err := s.mxResolver.Resolve(s.Ctx, fmt.Sprintf("%s.com", randomdata.Alphanumeric(32)))
	s.Nil(sl)
	s.NotNil(err)

	sl, err = s.mxResolver.Resolve(s.Ctx, "gmail.com")
	s.NotEmpty(sl)
	s.Equal(sl.Get(0).Host, "google.com")
	s.Nil(err)
}
