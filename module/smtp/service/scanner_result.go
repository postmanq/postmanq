package service

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"sync"
)

type ScannerResultStatus int

const (
	ScannerResultStatusProgress ScannerResultStatus = iota
	ScannerResultStatusSuccess
	ScannerResultStatusFailureMX
	ScannerResultStatusFailureIP
	ScannerResultStatusFailureIPLen
)

type ScannerResult interface {
	GetHostname() string
	GetStatus() ScannerResultStatus
	GetMxs() []entity.MX
	GetError() error
}

type result struct {
	hostname string
	status   ScannerResultStatus
	mxs      []entity.MX
	wg       *sync.WaitGroup
	err      error
}

func (s *result) GetHostname() string {
	return s.hostname
}

func (s *result) GetStatus() ScannerResultStatus {
	return s.status
}

func (s *result) GetMxs() []entity.MX {
	return s.mxs
}

func (s *result) lock() {
	s.wg.Add(1)
	s.status = ScannerResultStatusProgress
}

func (s *result) unlockWithStatus(status ScannerResultStatus) {
	s.status = status
	s.wg.Done()
}

func (s *result) GetError() error {
	return s.err
}
