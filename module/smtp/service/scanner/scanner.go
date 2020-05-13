package scanner

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"net"
	"sync"
)

type Scanner interface {
	Scan(string) Result
}

func NewScanner() Scanner {
	s := &scanner{
		mtx:           new(sync.Mutex),
		results:       make(map[string]*result),
		futureResults: make(chan *result, 1024),
	}
	go s.processFutureResults()
	return s
}

type scanner struct {
	mtx           *sync.Mutex
	results       map[string]*result
	futureResults chan *result
}

func (s *scanner) processFutureResults() {
	for futureResult := range s.futureResults {
		mxs, err := net.LookupMX(futureResult.hostname)
		if err != nil {
			//return err
		}

		futureResult.mxs = make([]entity.MX, len(mxs))
		for i, mx := range mxs {
			ips, err := net.LookupIP(mx.Host)
			if err != nil {

			}

			futureResult.mxs[i] = entity.MX{
				MX: mx,
				IP: ips[0],
			}

			futureResult.unlockWithStatus(ResultStatusSuccess)
		}
	}
}

func (s *scanner) Scan(hostname string) Result {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	r, ok := s.results[hostname]
	if !ok {
		r := &result{
			hostname: hostname,
			wg:       new(sync.WaitGroup),
		}

		r.lock()
		s.results[hostname] = r
		s.futureResults <- r

	}

	defer r.wg.Wait()
	return r
}
