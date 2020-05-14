package scanner_test

import (
	"github.com/postmanq/postmanq/module/smtp/service/scanner"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestScannerSuite(t *testing.T) {
	suite.Run(t, new(ScannerSuite))
}

type ScannerSuite struct {
	suite.Suite
	sc  scanner.Scanner
	mxs []string
	ips []string
}

func (s *ScannerSuite) SetupTest() {
	s.sc = scanner.NewScanner()
	s.mxs = []string{
		"aspmx.l.google.com",
		"alt1.aspmx.l.google.com",
		"alt2.aspmx.l.google.com",
		"alt3.aspmx.l.google.com",
		"alt4.aspmx.l.google.com",
	}
	s.ips = []string{
		"64.233.165.27",
		"64.233.189.27",
		"74.125.28.26",
		"108.177.8.27",
		"74.125.129.27",
	}
}

func (s *ScannerSuite) TestFailure() {
	result := s.sc.Scan("example.local")
	s.Equal(scanner.ResultStatusFailureMX, result.GetStatus())
}

func (s *ScannerSuite) TestSuccess() {
	result := s.sc.Scan("google.com")
	s.Equal(scanner.ResultStatusSuccess, result.GetStatus())
	s.Equal("google.com", result.GetHostname())
	for i, mx := range result.GetMxs() {
		s.Equal(s.mxs[i], mx.Host)
		s.Equal(s.ips[i], mx.IP.String())
	}
}
