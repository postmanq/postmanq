package service

import "net/smtp"

type ClientPool interface {
	GetFree() (*smtp.Client, error)
}

type clientPool struct {
}

func (p *clientPool) GetFree() (*smtp.Client, error) {
	return &smtp.Client{}, nil
}

type client struct {
}
