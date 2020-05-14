package scanner

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"sync"
)

type ResultStatus int

const (
	ResultStatusProgress ResultStatus = iota
	ResultStatusSuccess
	ResultStatusFailureMX
	ResultStatusFailureIP
	ResultStatusFailureIPLen
)

type Result interface {
	GetHostname() string
	GetStatus() ResultStatus
	GetMxs() []entity.MX
	GetError() error
}

type result struct {
	hostname string
	status   ResultStatus
	mxs      []entity.MX
	wg       *sync.WaitGroup
	err      error
}

func (s *result) GetHostname() string {
	return s.hostname
}

func (s *result) GetStatus() ResultStatus {
	return s.status
}

func (s *result) GetMxs() []entity.MX {
	return s.mxs
}

func (s *result) lock() {
	s.wg.Add(1)
	s.status = ResultStatusProgress
}

func (s *result) unlockWithStatus(status ResultStatus) {
	s.status = status
	s.wg.Done()
}

func (s *result) GetError() error {
	return s.err
}
