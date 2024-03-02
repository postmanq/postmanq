package services

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"os"
)

func NewFxDkimSignerFactory() smtp.DkimSignerFactory {
	return &dkimSignerFactory{}
}

type dkimSignerFactory struct {
}

func (d dkimSignerFactory) Create(ctx context.Context, cfg smtp.Config) (smtp.DkimSigner, error) {
	buf, err := os.ReadFile(cfg.TLS.PrivateKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, errors.New(fmt.Sprintf("could not decode PEM block from file %s", cfg.TLS.PrivateKey))
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &dkimSigner{
		opts: &dkim.SignOptions{
			Domain:   cfg.Hostname,
			Selector: cfg.DkimSelector,
			Signer:   privateKey,
		},
	}, nil
}

type dkimSigner struct {
	opts *dkim.SignOptions
}

func (d dkimSigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	var w bytes.Buffer
	r := bytes.NewReader(data)

	err := dkim.Sign(&w, r, d.opts)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
