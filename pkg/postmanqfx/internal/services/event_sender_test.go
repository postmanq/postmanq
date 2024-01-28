package services_test

import (
	"errors"
	"github.com/Pallinder/go-randomdata"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log_mock"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal_mocks"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/postmanqfx/internal/services"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq_mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	ErrActivity1 = errors.New("activity error 1")
	ErrActivity2 = errors.New("activity error 2")
	ErrWorkflow  = errors.New("workflow error")
)

func TestEventSenderTestSuite(t *testing.T) {
	suite.Run(t, new(EventSenderTestSuite))
}

type EventSenderTestSuite struct {
	testutils.Suite
	sender                  postmanq.EventSender
	workflowExecutorFactory *temporal_mocks.MockWorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	activityExecutorFactory *temporal_mocks.MockActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	middlewarePlugin        *postmanq_mocks.MockWorkflowPlugin
	senderPlugin            *postmanq_mocks.MockWorkflowPlugin
}

func (s *EventSenderTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.workflowExecutorFactory = temporal_mocks.NewMockWorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event](s.Ctrl)
	s.activityExecutorFactory = temporal_mocks.NewMockActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event](s.Ctrl)
	logger := log_mock.NewMockLogger(s.Ctrl)
	logger.EXPECT().Error(gomock.Any()).AnyTimes().Return()
	factory := services.NewFxEventSenderFactory(services.EventSenderFactoryParams{
		Ctx:                     s.Ctx,
		Logger:                  logger,
		WorkflowExecutorFactory: s.workflowExecutorFactory,
		ActivityExecutorFactory: s.activityExecutorFactory,
	})

	s.middlewarePlugin = postmanq_mocks.NewMockWorkflowPlugin(s.Ctrl)
	s.senderPlugin = postmanq_mocks.NewMockWorkflowPlugin(s.Ctrl)
	s.sender = factory.Create(&postmanq.Pipeline{
		Name:        randomdata.Title(randomdata.RandomGender),
		Middlewares: collection.ImportSlice[postmanq.WorkflowPlugin](s.middlewarePlugin),
		Senders:     collection.ImportSlice[postmanq.WorkflowPlugin](s.senderPlugin),
	})
}

func (s *EventSenderTestSuite) TestSendEvent() {
	middlewarePlugingType := randomdata.Alphanumeric(32)
	s.middlewarePlugin.EXPECT().GetType().AnyTimes().Return(middlewarePlugingType)
	middlewareActivityExecutor := temporal_mocks.NewMockActivityExecutor[*postmanqv1.Event, *postmanqv1.Event](s.Ctrl)
	middlewareActivityExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, ErrActivity1)
	s.activityExecutorFactory.EXPECT().Create(middlewarePlugingType).AnyTimes().Return(middlewareActivityExecutor)
	workflowExecutor := temporal_mocks.NewMockWorkflowExecutor[*postmanqv1.Event, *postmanqv1.Event](s.Ctrl)
	workflowExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, ErrWorkflow)
	s.workflowExecutorFactory.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(workflowExecutor)
	event, err := s.sender.SendEvent(nil, &postmanqv1.Event{})
	s.Nil(event)
	s.NotNil(err)
	s.ErrorIs(err, ErrActivity1)
	s.NotErrorIs(err, ErrActivity2)
	s.ErrorIs(err, ErrWorkflow)

	middlewareActivityExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, ErrActivity1)
	workflowExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
	event, err = s.sender.SendEvent(nil, &postmanqv1.Event{})
	s.Nil(event)
	s.NotNil(err)
	s.ErrorIs(err, ErrActivity1)
	s.NotErrorIs(err, ErrActivity2)
	s.NotErrorIs(err, ErrWorkflow)

	senderPluginType := randomdata.Alphanumeric(32)
	s.senderPlugin.EXPECT().GetType().AnyTimes().Return(senderPluginType)
	senderActivityExecutor := temporal_mocks.NewMockActivityExecutor[*postmanqv1.Event, *postmanqv1.Event](s.Ctrl)
	s.activityExecutorFactory.EXPECT().Create(senderPluginType).AnyTimes().Return(senderActivityExecutor)

	middlewareActivityExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
	senderActivityExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, ErrActivity2)
	event, err = s.sender.SendEvent(nil, &postmanqv1.Event{})
	s.Nil(event)
	s.NotNil(err)
	s.ErrorIs(err, ErrActivity2)
	s.NotErrorIs(err, ErrActivity1)
	s.NotErrorIs(err, ErrWorkflow)

	senderActivityExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, nil)
	event, err = s.sender.SendEvent(nil, &postmanqv1.Event{})
	s.NotNil(event)
	s.Nil(err)
}
