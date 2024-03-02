package services_test

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/internal/services"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

const (
	privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIICXwIBAAKBgQDwIRP/UC3SBsEmGqZ9ZJW3/DkMoGeLnQg1fWn7/zYtIxN2SnFC
jxOCKG9v3b4jYfcTNh5ijSsq631uBItLa7od+v/RtdC2UzJ1lWT947qR+Rcac2gb
to/NMqJ0fzfVjH4OuKhitdY9tf6mcwGjaNBcWToIMmPSPDdQPNUYckcQ2QIDAQAB
AoGBALmn+XwWk7akvkUlqb+dOxyLB9i5VBVfje89Teolwc9YJT36BGN/l4e0l6QX
/1//6DWUTB3KI6wFcm7TWJcxbS0tcKZX7FsJvUz1SbQnkS54DJck1EZO/BLa5ckJ
gAYIaqlA9C0ZwM6i58lLlPadX/rtHb7pWzeNcZHjKrjM461ZAkEA+itss2nRlmyO
n1/5yDyCluST4dQfO8kAB3toSEVc7DeFeDhnC1mZdjASZNvdHS4gbLIA1hUGEF9m
3hKsGUMMPwJBAPW5v/U+AWTADFCS22t72NUurgzeAbzb1HWMqO4y4+9Hpjk5wvL/
eVYizyuce3/fGke7aRYw/ADKygMJdW8H/OcCQQDz5OQb4j2QDpPZc0Nc4QlbvMsj
7p7otWRO5xRa6SzXqqV3+F0VpqvDmshEBkoCydaYwc2o6WQ5EBmExeV8124XAkEA
qZzGsIxVP+sEVRWZmW6KNFSdVUpk3qzK0Tz/WjQMe5z0UunY9Ax9/4PVhp/j61bf
eAYXunajbBSOLlx4D+TunwJBANkPI5S9iylsbLs6NkaMHV6k5ioHBBmgCak95JGX
GMot/L2x0IYyMLAz6oLWh2hm7zwtb0CgOrPo1ke44hFYnfc=
-----END RSA PRIVATE KEY-----
`
	mailData = `From: Test One <test-1@postmanq.io>
To: Test Two <test-2@postmanq.io>
Subject: Test message
Date: Fri, 01 Feb 2024 00:00:0 -0000 (PDT)
Message-ID: <c586127a-7fcc-4401-b538-1dbfe87bcec4@postmanq.io>

Hello world!
`
)

func TestDkimTestSuite(t *testing.T) {
	suite.Run(t, new(DkimTestSuite))
}

type DkimTestSuite struct {
	testutils.Suite
	privateKeyFile *os.File
	factory        smtp.DkimSignerFactory
}

func (s *DkimTestSuite) SetupTest() {
	privateKeyFile, err := os.CreateTemp("/tmp", "postmanq_private_key_*.pem")
	s.Nil(err)
	_, err = privateKeyFile.Write([]byte(privateKeyPEM))
	s.Nil(err)

	s.privateKeyFile = privateKeyFile
	s.factory = services.NewFxDkimSignerFactory()
}

func (s *DkimTestSuite) TestCreateSigner() {
	_, err := s.factory.Create(s.Ctx, smtp.Config{
		TLS: &smtp.TLSConfig{
			PrivateKey: randomdata.Alphanumeric(32),
		},
	})
	s.NotNil(err)

	_, err = s.privateKeyFile.WriteAt([]byte(randomdata.Alphanumeric(100)), 10)
	s.Nil(err)

	_, err = s.factory.Create(s.Ctx, smtp.Config{
		TLS: &smtp.TLSConfig{
			PrivateKey: s.privateKeyFile.Name(),
		},
	})
	s.NotNil(err)
}

func (s *DkimTestSuite) TestSign() {
	signer, err := s.factory.Create(s.Ctx, smtp.Config{
		Hostname:     "postmanq.io",
		DkimSelector: "_dkim",
		TLS: &smtp.TLSConfig{
			PrivateKey: s.privateKeyFile.Name(),
		},
	})
	s.NotNil(signer)
	s.Nil(err)

	_, err = signer.Sign(s.Ctx, []byte(""))
	s.NotNil(err)

	_, err = signer.Sign(s.Ctx, []byte(mailData))
	s.Nil(err)
}

func (s *DkimTestSuite) TearDownTest() {
	err := os.Remove(s.privateKeyFile.Name())
	s.Nil(err)
}
