package service

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"net/smtp"
)

type connector struct {
}

func (c *connector) Connect(mx entity.MX) (*smtp.Client, error) {
	return &smtp.Client{}, nil
}
