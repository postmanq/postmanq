package component_test

import (
	"errors"
	cs "github.com/postmanq/postmanq/mock/module/config/service"
	rs "github.com/postmanq/postmanq/mock/module/rabbitmq/service"
	vs "github.com/postmanq/postmanq/mock/module/validator/service"
	"github.com/postmanq/postmanq/module/rabbitmq/component"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	defaultError     = errors.New("")
	expectedCallsNum = 1
)

func TestReceiverSuite(t *testing.T) {
	suite.Run(t, new(ReceiverSuite))
}

type ReceiverSuite struct {
	suite.Suite
	configProvider *cs.ConfigProvider
	pool           *rs.Pool
	validator      *vs.Validator
	receiver       component.Receiver
}

func (s *ReceiverSuite) SetupTest() {
	s.configProvider = new(cs.ConfigProvider)
	s.pool = new(rs.Pool)
	s.validator = new(vs.Validator)
	s.receiver = component.NewReceiver(
		s.configProvider,
		s.pool,
		s.validator,
	).Descriptor.Construct(nil).(component.Receiver)
}

func (s *ReceiverSuite) TestFailureConfig() {
	s.configProvider.On("Populate", mock.Anything, mock.Anything).Return(defaultError)
	s.NotNil(s.receiver.OnInit())
	s.configProvider.AssertNumberOfCalls(s.T(), "Populate", expectedCallsNum)
}

func (s *ReceiverSuite) TestFailureValidator() {
	s.configProvider.On("Populate", mock.Anything, mock.Anything).Return(nil)
	s.validator.On("Struct", mock.Anything).Return(defaultError)
	s.NotNil(s.receiver.OnInit())
	s.configProvider.AssertNumberOfCalls(s.T(), "Populate", expectedCallsNum)
	s.validator.AssertNumberOfCalls(s.T(), "Struct", expectedCallsNum)
}

func (s *ReceiverSuite) TestFailureConnect() {
	s.configProvider.On("Populate", mock.Anything, mock.Anything).Return(nil)
	s.validator.On("Struct", mock.Anything).Return(nil)
	s.pool.On("Connect", mock.Anything).Return(defaultError)
	s.NotNil(s.receiver.OnInit())
	s.configProvider.AssertNumberOfCalls(s.T(), "Populate", expectedCallsNum)
	s.validator.AssertNumberOfCalls(s.T(), "Struct", expectedCallsNum)
	s.pool.AssertNumberOfCalls(s.T(), "Connect", expectedCallsNum)
}

func (s *ReceiverSuite) TestFailureCreateSubscriber() {
	s.configProvider.On("Populate", mock.Anything, mock.Anything).Return(nil)
	s.validator.On("Struct", mock.Anything).Return(nil)
	s.pool.On("Connect", mock.Anything).Return(nil)
	s.pool.On("CreateSubscriber", mock.Anything, mock.Anything).Return(nil, defaultError)
	s.NotNil(s.receiver.OnInit())
	s.configProvider.AssertNumberOfCalls(s.T(), "Populate", expectedCallsNum)
	s.validator.AssertNumberOfCalls(s.T(), "Struct", expectedCallsNum)
	s.pool.AssertNumberOfCalls(s.T(), "Connect", expectedCallsNum)
	s.pool.AssertNumberOfCalls(s.T(), "CreateSubscriber", expectedCallsNum)
}

func (s *ReceiverSuite) TestSuccessOnInit() {
	s.configProvider.On("Populate", mock.Anything, mock.Anything).Return(nil)
	s.validator.On("Struct", mock.Anything).Return(nil)
	s.pool.On("Connect", mock.Anything).Return(nil)
	s.pool.On("CreateSubscriber", mock.Anything, mock.Anything).Return(nil, nil)
	s.Nil(s.receiver.OnInit())
	s.configProvider.AssertNumberOfCalls(s.T(), "Populate", expectedCallsNum)
	s.validator.AssertNumberOfCalls(s.T(), "Struct", expectedCallsNum)
	s.pool.AssertNumberOfCalls(s.T(), "Connect", expectedCallsNum)
	s.pool.AssertNumberOfCalls(s.T(), "CreateSubscriber", expectedCallsNum)
}
