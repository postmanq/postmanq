package scanner

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"sync"
)

type ResultStatus int

const (
	ResultStatusProgress ResultStatus = iota
	ResultStatusSuccess
	ResultStatusFailure
)

type Result interface {
	GetHostname() string
	GetStatus() ResultStatus
	GetMxs() []entity.MX
}

type result struct {
	hostname string
	status   ResultStatus
	mxs      []entity.MX
	wg       *sync.WaitGroup
}

func (s *result) GetHostname() string {
	defer s.wg.Wait()
	return s.hostname
}

func (s *result) GetStatus() ResultStatus {
	defer s.wg.Wait()
	return s.status
}

func (s *result) GetMxs() []entity.MX {
	defer s.wg.Wait()
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
