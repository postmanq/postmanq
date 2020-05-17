package service

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"github.com/postmanq/postmanq/module/smtp/errors"
	"net"
	"strings"
	"sync"
	"time"
)

type Scanner interface {
	Scan(string) ScannerResult
}

func NewScanner() Scanner {
	s := &scanner{
		scanMtx:            new(sync.Mutex),
		rescanMtx:          new(sync.Mutex),
		results:            make(map[string]*result),
		processableResults: make(chan *result, 1024),
		failureResults:     make(map[string][]*result),
		ticker:             time.NewTicker(time.Minute),
		layout:             "15:04",
	}
	go s.processResults()
	go s.rescanResults()
	go s.rescanFailureResults()
	return s
}

type scanner struct {
	scanMtx            *sync.Mutex
	rescanMtx          *sync.Mutex
	results            map[string]*result
	failureResults     map[string][]*result
	ticker             *time.Ticker
	processableResults chan *result
	layout             string
}

func (s *scanner) processResults() {
	for processableResult := range s.processableResults {
		status, err := s.processResult(processableResult)
		if err != nil {
			s.rescanMtx.Lock()
			key := time.Now().Add(time.Minute * 10).Format(s.layout)
			if _, ok := s.failureResults[key]; !ok {
				s.failureResults[key] = make([]*result, 0)
			}

			s.failureResults[key] = append(s.failureResults[key], processableResult)
			s.rescanMtx.Unlock()
		}

		processableResult.err = err
		processableResult.unlockWithStatus(status)
	}
}

func (s *scanner) processResult(processableResult *result) (ScannerResultStatus, error) {
	mxs, err := net.LookupMX(processableResult.hostname)
	if err != nil {
		return ScannerResultStatusFailureMX, err
	}

	processableResult.mxs = make([]entity.MX, len(mxs))
	for i, mx := range mxs {
		ips, err := net.LookupIP(mx.Host)
		if err != nil {
			return ScannerResultStatusFailureIP, err
		}

		if len(ips) == 0 {
			return ScannerResultStatusFailureIPLen, errors.IPsIsNotFoundByMX(mx.Host)
		}

		mx.Host = strings.TrimRight(mx.Host, ".")
		processableResult.mxs[i] = entity.MX{
			MX: mx,
			IP: ips[0],
		}
	}

	return ScannerResultStatusSuccess, nil
}

func (s *scanner) Scan(hostname string) ScannerResult {
	s.scanMtx.Lock()
	defer s.scanMtx.Unlock()
	r, ok := s.results[hostname]
	if !ok {
		r = &result{
			hostname: hostname,
			wg:       new(sync.WaitGroup),
		}

		r.lock()
		s.results[hostname] = r
		s.processableResults <- r
	}

	defer r.wg.Wait()
	return r
}

func (s *scanner) rescanResults() {
	for now := range s.ticker.C {
		if now.Minute() == 0 {
			s.scanMtx.Lock()
			for _, result := range s.results {
				if result.status == ScannerResultStatusSuccess {
					result.lock()
					s.processableResults <- result
				}
			}
			s.scanMtx.Unlock()
		}
	}
}

func (s *scanner) rescanFailureResults() {
	for now := range s.ticker.C {
		s.rescanMtx.Lock()
		key := now.Format(s.layout)
		if results, ok := s.failureResults[key]; ok {
			for _, result := range results {
				result.lock()
				s.processableResults <- result
			}

			delete(s.failureResults, key)
		}

		s.rescanMtx.Unlock()
	}
}
