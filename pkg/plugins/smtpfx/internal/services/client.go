package services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/thegeekyasian/round-robin-go"
	"net"
	smtp2 "net/smtp"
	"os"
	"sync"
	"time"
)

func NewFxClientBuilderFactory(
	logger log.Logger,
) smtp.ClientBuilderFactory {
	return &clientBuilderFactory{
		logger: logger,
	}
}

type clientBuilderFactory struct {
	logger log.Logger
}

func (c clientBuilderFactory) Create(ctx context.Context, cfg smtp.Config) (smtp.ClientBuilder, error) {
	var tlsCfg *tls.Config

	if cfg.TLS != nil {
		tlsCfg = &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
			MinVersion:             tls.VersionTLS12,
			SessionTicketsDisabled: true,
		}

		buf, err := os.ReadFile(cfg.TLS.Certificate)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(buf)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		pool.AddCert(cert)
		tlsCfg.RootCAs = pool
		tlsCfg.ClientCAs = pool
		x509KeyPair, err := tls.LoadX509KeyPair(cfg.TLS.Certificate, cfg.TLS.PrivateKey)
		if err == nil {
			return nil, err
		}

		tlsCfg.Certificates = []tls.Certificate{x509KeyPair}
	}

	ips, err := roundrobin.New[string](cfg.IPs...)
	if err != nil {
		return nil, err
	}

	return &clientBuilder{
		logger:     c.logger,
		ips:        ips,
		timeoutCfg: cfg.Timeout,
		tlsCfg:     tlsCfg,
	}, nil
}

type clientBuilder struct {
	logger     log.Logger
	ips        *roundrobin.RoundRobin[string]
	timeoutCfg *smtp.TimeoutConfig
	tlsCfg     *tls.Config
}

func (c *clientBuilder) Create(ctx context.Context, hostname string) (smtp.Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(*c.ips.Next(), "0"))
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{
		Timeout:   c.timeoutCfg.Connection,
		LocalAddr: tcpAddr,
	}
	conn, err := dialer.Dial("tcp", net.JoinHostPort(hostname, "25"))
	if err != nil {
		return nil, err
	}

	smtpClient, err := smtp2.NewClient(conn, hostname)
	if err != nil {
		return nil, err
	}

	return &client{
		logger:     c.logger,
		conn:       conn,
		smtpClient: smtpClient,
		timeoutCfg: c.timeoutCfg,
		tlsCfg:     c.tlsCfg,
		status:     smtp.ClientStatusBusy,
		tlsStatus:  clientTLSStatusUndefined,
	}, nil
}

type clientTLSStatus int

const (
	clientTLSStatusUndefined clientTLSStatus = iota
	clientTLSStatusUse
	clientTLSStatusDontUse
)

type client struct {
	logger     log.Logger
	conn       net.Conn
	smtpClient *smtp2.Client
	timeoutCfg *smtp.TimeoutConfig
	tlsCfg     *tls.Config
	status     smtp.ClientStatus
	tlsStatus  clientTLSStatus
	mtx        sync.Mutex
}

func (c *client) Hello(ctx context.Context, localName string) error {
	err := c.conn.SetDeadline(time.Now().Add(c.timeoutCfg.Hello))
	if err != nil {
		return err
	}

	err = c.smtpClient.Hello(localName)
	if err != nil {
		return err
	}

	if c.tlsStatus == clientTLSStatusUndefined {
		useTLS, _ := c.smtpClient.Extension("STARTTLS")
		if useTLS {
			c.tlsStatus = clientTLSStatusUse
		} else {
			c.tlsStatus = clientTLSStatusDontUse
		}
	}

	if c.tlsStatus == clientTLSStatusUse {
		return c.smtpClient.StartTLS(c.tlsCfg)
	}

	return nil
}

func (c *client) Mail(ctx context.Context, from string) error {
	err := c.conn.SetDeadline(time.Now().Add(c.timeoutCfg.Mail))
	if err != nil {
		return err
	}

	return c.smtpClient.Mail(from)
}

func (c *client) Rcpt(ctx context.Context, to string) error {
	err := c.conn.SetDeadline(time.Now().Add(c.timeoutCfg.Rcpt))
	if err != nil {
		return err
	}

	return c.smtpClient.Rcpt(to)
}

func (c *client) Data(ctx context.Context, data []byte) error {
	err := c.conn.SetDeadline(time.Now().Add(c.timeoutCfg.Data))
	if err != nil {
		return err
	}

	wc, err := c.smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write(data)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	err = c.smtpClient.Reset()
	if err != nil {
		return err
	}

	c.mtx.Lock()
	c.status = smtp.ClientStatusFree
	c.mtx.Unlock()
	return nil
}

func (c *client) HasStatus(expectedStatus smtp.ClientStatus) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.status == expectedStatus
}

func (c *client) Noop() error {
	return c.smtpClient.Noop()
}
