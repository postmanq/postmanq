package service

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"net"
	"sync"
)

type Scanner interface {
	Get() *entity.MailServer
	Scan(string) error
	//ScanByStatus(entity.MXServerStatus) error
}

func NewScanner() Scanner {
	return &scanner{
		//mxs: make(entity.MXServers),
		//mtx: new(sync.Mutex),
	}
}

type scanner struct {
	//mxs entity.MXServers
	mtx *sync.Mutex
}

func (s *scanner) Get() *entity.MailServer {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	//mxServers, ok := s.mxs[hostname]
	//if !ok {
	//	return nil, errors.MXServersIsNotFound(hostname)
	//}

	return nil
}

func (s *scanner) Scan(hostname string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	mxs, err := net.LookupMX(hostname)
	if err != nil {
		return err
	}

	mxServers := make([]entity.MX, len(mxs))
	for i, mx := range mxs {
		ip, err := net.LookupIP(mx.Host)

	}

	//s.mxs[hostname] = mxServers
	return nil
}
