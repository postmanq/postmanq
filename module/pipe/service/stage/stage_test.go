package stage_test

import (
	"fmt"
	mm "github.com/postmanq/postmanq/mock/module"
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/service/stage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/multierr"
	"testing"
)

func TestStageSuite(t *testing.T) {
	suite.Run(t, new(StageSuite))
}

type StageSuite struct {
	suite.Suite
	receiveStage            stage.ResultStage
	middlewareStage         stage.DeliveryStage
	parallelMiddlewareStage stage.DeliveryStage
	completeStage           stage.DeliveryStage
	receiver                *mm.ReceiveComponent
	middleware              *mm.ProcessComponent
	parallelMiddleware1     *mm.ProcessComponent
	parallelMiddleware2     *mm.ProcessComponent
	parallelMiddleware3     *mm.ProcessComponent
	sender                  *mm.SendComponent
}

func (s *StageSuite) SetupTest() {
	s.receiver = new(mm.ReceiveComponent)
	s.sender = new(mm.SendComponent)
	s.middleware = new(mm.ProcessComponent)
	s.parallelMiddleware1 = new(mm.ProcessComponent)
	s.parallelMiddleware2 = new(mm.ProcessComponent)
	s.parallelMiddleware3 = new(mm.ProcessComponent)

	receiveStage, err := stage.NewReceive().Constructor(nil, s.receiver)
	s.Nil(err)
	s.receiveStage = receiveStage.(stage.ResultStage)

	completeStage, err := stage.NewComplete().Constructor(nil, s.sender)
	s.Nil(err)
	s.completeStage = completeStage.(stage.DeliveryStage)

	middlewareStage, err := stage.NewMiddleware().Constructor(nil, s.middleware)
	s.Nil(err)
	s.middlewareStage = middlewareStage.(stage.DeliveryStage)

	parallelMiddlewareStage, err := stage.NewParallelMiddleware().Constructor(nil, []interface{}{
		s.parallelMiddleware1,
		s.parallelMiddleware2,
		s.parallelMiddleware3,
	})
	s.Nil(err)
	s.parallelMiddlewareStage = parallelMiddlewareStage.(stage.DeliveryStage)

	s.receiveStage.Bind(s.middlewareStage)

	s.middlewareStage.Bind(s.receiveStage)
	s.middlewareStage.Bind(s.parallelMiddlewareStage)

	s.parallelMiddlewareStage.Bind(s.receiveStage)
	s.parallelMiddlewareStage.Bind(s.completeStage)

	s.completeStage.Bind(s.receiveStage)
}

func (s *StageSuite) TestReceiveStageFailure() {
	s.receiver.On("OnReceive", mock.Anything, mock.Anything).Return(fmt.Errorf("receiver error")).Once()
	s.NotNil(s.receiveStage.Start())
	s.receiver.AssertNumberOfCalls(s.T(), "OnReceive", 1)
}

func (s *StageSuite) TestReceiveStageSuccess() {
	s.receiver.On("OnReceive", mock.Anything, mock.Anything).Return(nil).Once()
	s.Nil(s.receiveStage.Start())
	s.receiver.AssertNumberOfCalls(s.T(), "OnReceive", 1)
}

func (s *StageSuite) TestMiddlewareStage() {
	err := fmt.Errorf("middleware error")
	go func() {
		err := s.middlewareStage.Start()
		s.Nil(err)
	}()

	s.middleware.On("OnProcess", mock.Anything).Return(err).Once()
	s.middlewareStage.Deliveries() <- module.Delivery{}
	result := <-s.receiveStage.Results()
	s.NotNil(result)
	s.Equal(err, result.Err)

	s.middleware.On("OnProcess", mock.Anything).Return(nil).Once()
	s.middlewareStage.Deliveries() <- module.Delivery{}
	delivery := <-s.parallelMiddlewareStage.Deliveries()
	s.NotNil(delivery)
	s.Nil(delivery.Err)
	s.middleware.AssertNumberOfCalls(s.T(), "OnProcess", 2)
}

func (s *StageSuite) TestParallelMiddlewareStage() {
	err := fmt.Errorf("parallel middleware error")
	combinedErr := multierr.Combine(err, err)
	go func() {
		err := s.parallelMiddlewareStage.Start()
		s.Nil(err)
	}()

	s.parallelMiddleware1.On("OnProcess", mock.Anything).Return(err).Once()
	s.parallelMiddleware2.On("OnProcess", mock.Anything).Return(err).Once()
	s.parallelMiddleware3.On("OnProcess", mock.Anything).Return(nil).Once()
	s.parallelMiddlewareStage.Deliveries() <- module.Delivery{}
	result := <-s.receiveStage.Results()
	s.NotNil(result)
	s.Equal(combinedErr, result.Err)

	s.parallelMiddleware1.On("OnProcess", mock.Anything).Return(nil).Once()
	s.parallelMiddleware2.On("OnProcess", mock.Anything).Return(nil).Once()
	s.parallelMiddleware3.On("OnProcess", mock.Anything).Return(nil).Once()
	s.parallelMiddlewareStage.Deliveries() <- module.Delivery{}
	delivery := <-s.completeStage.Deliveries()
	s.NotNil(delivery)
	s.Nil(delivery.Err)
	s.parallelMiddleware1.AssertNumberOfCalls(s.T(), "OnProcess", 2)
	s.parallelMiddleware2.AssertNumberOfCalls(s.T(), "OnProcess", 2)
	s.parallelMiddleware3.AssertNumberOfCalls(s.T(), "OnProcess", 2)
}

func (s *StageSuite) TestCompleteStage() {
	err := fmt.Errorf("complete error")
	go func() {
		err := s.completeStage.Start()
		s.Nil(err)
	}()

	s.sender.On("OnSend", mock.Anything).Return(err).Once()
	s.completeStage.Deliveries() <- module.Delivery{}
	result := <-s.receiveStage.Results()
	s.NotNil(result)
	s.Equal(err, result.Err)

	s.sender.On("OnSend", mock.Anything).Return(nil).Once()
	s.completeStage.Deliveries() <- module.Delivery{}
	result = <-s.receiveStage.Results()
	s.NotNil(result)
	s.Nil(result.Err)
}
