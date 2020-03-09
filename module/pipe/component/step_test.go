package component_test

import (
	"errors"
	"github.com/byorty/postmanq_sender/plugin"
	"github.com/byorty/postmanq_sender/sending"
	"github.com/byorty/postmanq_sender/util/chain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/multierr"
	"runtime"
	"sync"
	"testing"
)

type MockSender struct {
	mock.Mock
}

func (m *MockSender) OnSend(s sending.Sending) error {
	args := m.Called(s)
	return args.Error(0)
}

type MockSending struct {
	mock.Mock
	wg *sync.WaitGroup
}

func (m *MockSending) Complete() {
	m.Called()
	m.wg.Done()
}

func (m *MockSending) Abort(err error) {
	m.Called(err)
	m.wg.Done()
}

func (m *MockSending) Wait() {
	m.wg.Wait()
}

type MockSendingChain struct {
	mock.Mock
}

func (m *MockSendingChain) Run(numCPU int) {}

func (m *MockSendingChain) Send(s sending.Sending) {
	m.Called(s)
	s.(*MockSending).Complete()
}

func TestChain(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	mockSender         *MockSender
	mockNextChain      *MockSendingChain
	mockSending        *MockSending
	numCPU             int
	simpleNumOfCalls   int
	simpleErr          error
	parallelNumOfCalls int
	parallelErr        error
}

func (s *Suite) SetupTest() {
	s.numCPU = runtime.NumCPU()
	s.simpleNumOfCalls = 1
	s.simpleErr = errors.New("sender error")
	s.parallelNumOfCalls = 3
	s.parallelErr = s.simpleErr
	s.parallelErr = multierr.Append(s.parallelErr, s.simpleErr)
	s.parallelErr = multierr.Append(s.parallelErr, s.simpleErr)

	s.mockSender = new(MockSender)

	s.mockSending = new(MockSending)
	s.mockSending.wg = new(sync.WaitGroup)
	s.mockSending.wg.Add(1)
	s.mockSending.On("Complete").Return(nil)

	s.mockNextChain = new(MockSendingChain)
	s.mockNextChain.On("Send", s.mockSending).Return(nil)
}

func (s *Suite) TestSimpleError() {
	s.mockSender.On("OnSend", s.mockSending).Return(s.simpleErr)
	s.mockSending.On("Abort", s.simpleErr).Return(nil)

	simple := chain.NewSimpleSender(s.mockSender, s.mockNextChain)
	simple.Run(s.numCPU)
	simple.Send(s.mockSending)

	s.mockSending.Wait()

	s.mockSender.AssertCalled(s.T(), "OnSend", s.mockSending)
	s.mockSender.AssertNumberOfCalls(s.T(), "OnSend", s.simpleNumOfCalls)
	s.mockSending.AssertCalled(s.T(), "Abort", s.simpleErr)
	s.mockSending.AssertNotCalled(s.T(), "Complete")
	s.mockNextChain.AssertNotCalled(s.T(), "Send")
}

func (s *Suite) TestSimpleSuccess() {
	s.mockSender.On("OnSend", s.mockSending).Return(nil)

	simple := chain.NewSimpleSender(s.mockSender, s.mockNextChain)
	simple.Run(s.numCPU)
	simple.Send(s.mockSending)

	s.mockSending.Wait()

	s.mockSender.AssertCalled(s.T(), "OnSend", s.mockSending)
	s.mockSender.AssertNumberOfCalls(s.T(), "OnSend", s.simpleNumOfCalls)
	s.mockNextChain.AssertCalled(s.T(), "Send", s.mockSending)
	s.mockNextChain.AssertNumberOfCalls(s.T(), "Send", s.simpleNumOfCalls)
}

func (s *Suite) TestCompleteError() {
	s.mockSender.On("OnSend", s.mockSending).Return(s.simpleErr)
	s.mockSending.On("Abort", s.simpleErr).Return(nil)

	complete := chain.NewCompleteSender(s.mockSender)
	complete.Run(s.numCPU)
	complete.Send(s.mockSending)

	s.mockSending.Wait()

	s.mockSender.AssertCalled(s.T(), "OnSend", s.mockSending)
	s.mockSender.AssertNumberOfCalls(s.T(), "OnSend", s.simpleNumOfCalls)
	s.mockSending.AssertCalled(s.T(), "Abort", s.simpleErr)
	s.mockSending.AssertNotCalled(s.T(), "Complete")
}

func (s *Suite) TestCompleteSuccess() {
	s.mockSender.On("OnSend", s.mockSending).Return(nil)

	complete := chain.NewCompleteSender(s.mockSender)
	complete.Run(s.numCPU)
	complete.Send(s.mockSending)

	s.mockSending.Wait()

	s.mockSender.AssertCalled(s.T(), "OnSend", s.mockSending)
	s.mockSender.AssertNumberOfCalls(s.T(), "OnSend", s.simpleNumOfCalls)
	s.mockSending.AssertCalled(s.T(), "Complete")
}

func (s *Suite) TestParallelError() {
	s.mockSender.On("OnSend", s.mockSending).Return(s.simpleErr)
	s.mockSending.On("Abort", s.parallelErr).Return(nil)

	parallel := chain.NewParallelSender([]plugin.Sender{s.mockSender, s.mockSender, s.mockSender}, s.mockNextChain)
	parallel.Run(s.numCPU)
	parallel.Send(s.mockSending)

	s.mockSending.Wait()

	s.mockSender.AssertCalled(s.T(), "OnSend", s.mockSending)
	s.mockSender.AssertNumberOfCalls(s.T(), "OnSend", s.parallelNumOfCalls)
	s.mockSending.AssertCalled(s.T(), "Abort", s.parallelErr)
	s.mockSending.AssertNotCalled(s.T(), "Complete")
	s.mockNextChain.AssertNotCalled(s.T(), "Send")
}

func (s *Suite) TestParallelSuccess() {
	s.mockSender.On("OnSend", s.mockSending).Return(nil)

	parallel := chain.NewParallelSender([]plugin.Sender{s.mockSender, s.mockSender, s.mockSender}, s.mockNextChain)
	parallel.Run(s.numCPU)
	parallel.Send(s.mockSending)

	s.mockSending.Wait()

	s.mockSender.AssertCalled(s.T(), "OnSend", s.mockSending)
	s.mockSender.AssertNumberOfCalls(s.T(), "OnSend", s.parallelNumOfCalls)
	s.mockNextChain.AssertCalled(s.T(), "Send", s.mockSending)
	s.mockNextChain.AssertNumberOfCalls(s.T(), "Send", s.simpleNumOfCalls)
}
