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
	"time"
)

type ReceiveComponent struct {
	out chan module.Delivery
}

func (c *ReceiveComponent) OnReceive(out chan module.Delivery) error {
	c.out = out
	return nil
}

func TestStageSuite(t *testing.T) {
	suite.Run(t, new(StageSuite))
}

type StageSuite struct {
	suite.Suite
	receiveStage            stage.Stage
	middlewareStage         stage.Stage
	parallelMiddlewareStage stage.Stage
	completeStage           stage.Stage
	receiver                *ReceiveComponent
	middleware              *mm.ProcessComponent
	parallelMiddleware1     *mm.ProcessComponent
	parallelMiddleware2     *mm.ProcessComponent
	parallelMiddleware3     *mm.ProcessComponent
	sender                  *mm.SendComponent
}

func (s *StageSuite) SetupTest() {
	s.receiver = new(ReceiveComponent)
	s.sender = new(mm.SendComponent)
	s.middleware = new(mm.ProcessComponent)
	s.parallelMiddleware1 = new(mm.ProcessComponent)
	s.parallelMiddleware2 = new(mm.ProcessComponent)
	s.parallelMiddleware3 = new(mm.ProcessComponent)

	receiveStage, err := stage.NewReceive().Constructor(nil, s.receiver)
	s.Nil(err)
	s.receiveStage = receiveStage

	completeStage, err := stage.NewComplete().Constructor(nil, s.sender)
	s.Nil(err)
	s.completeStage = completeStage

	middlewareStage, err := stage.NewMiddleware().Constructor(nil, s.middleware)
	s.Nil(err)
	s.middlewareStage = middlewareStage

	parallelMiddlewareStage, err := stage.NewParallelMiddleware().Constructor(nil, []interface{}{
		s.parallelMiddleware1,
		s.parallelMiddleware2,
		s.parallelMiddleware3,
	})
	s.Nil(err)
	s.parallelMiddlewareStage = parallelMiddlewareStage
}

func (s *StageSuite) TestReceiveStage() {
	receiveOut := s.receiveStage.Start(nil)
	time.Sleep(time.Second)
	s.receiver.out <- module.Delivery{}
	d := <-receiveOut
	s.NotNil(d)
}

func (s *StageSuite) TestMiddlewareStage() {
	receiveOut := s.receiveStage.Start(nil)
	err := fmt.Errorf("middleware error")
	middlewareOut := s.middlewareStage.Start(receiveOut)
	time.Sleep(time.Second)

	s.middleware.On("OnProcess", mock.Anything).Return(err).Once()
	d := module.Delivery{
		Err: make(chan error, 1),
	}
	s.receiver.out <- d
	e := <-d.Err
	s.NotNil(e)
	s.Equal(err, e)

	s.middleware.On("OnProcess", mock.Anything).Return(nil).Once()
	s.receiver.out <- d
	d = <-middlewareOut
	s.NotNil(d)
	s.middleware.AssertNumberOfCalls(s.T(), "OnProcess", 2)
}

func (s *StageSuite) TestParallelMiddlewareStage() {
	receiveOut := s.receiveStage.Start(nil)
	err := fmt.Errorf("parallel middleware error")
	var multiErr error
	multiErr = multierr.Append(multiErr, err)
	multiErr = multierr.Append(multiErr, err)
	parallelMiddlewareOut := s.parallelMiddlewareStage.Start(receiveOut)
	time.Sleep(time.Second)

	s.parallelMiddleware1.On("OnProcess", mock.Anything).Return(err).Once()
	s.parallelMiddleware2.On("OnProcess", mock.Anything).Return(err).Once()
	s.parallelMiddleware3.On("OnProcess", mock.Anything).Return(nil).Once()
	d := module.Delivery{
		Err: make(chan error, 1),
	}
	s.receiver.out <- d
	e := <-d.Err
	s.NotNil(e)
	s.Equal(multiErr, e)

	s.parallelMiddleware1.On("OnProcess", mock.Anything).Return(nil).Once()
	s.parallelMiddleware2.On("OnProcess", mock.Anything).Return(nil).Once()
	s.parallelMiddleware3.On("OnProcess", mock.Anything).Return(nil).Once()
	s.receiver.out <- d
	delivery := <-parallelMiddlewareOut
	s.NotNil(delivery)
	s.parallelMiddleware1.AssertNumberOfCalls(s.T(), "OnProcess", 2)
	s.parallelMiddleware2.AssertNumberOfCalls(s.T(), "OnProcess", 2)
	s.parallelMiddleware3.AssertNumberOfCalls(s.T(), "OnProcess", 2)
}

func (s *StageSuite) TestCompleteStage() {
	receiveOut := s.receiveStage.Start(nil)
	err := fmt.Errorf("complete error")
	s.completeStage.Start(receiveOut)
	time.Sleep(time.Second)

	s.sender.On("OnSend", mock.Anything).Return(err).Once()
	d := module.Delivery{
		Err: make(chan error, 1),
	}
	s.receiver.out <- d
	e := <-d.Err
	s.NotNil(e)
	s.Equal(err, e)

	s.sender.On("OnSend", mock.Anything).Return(nil).Once()
	d = module.Delivery{
		Err: make(chan error, 1),
	}
	s.receiver.out <- d
	s.Nil(<-d.Err)
}
