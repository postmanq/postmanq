package service

import (
	"github.com/postmanq/postmanq/module/smtp/entity"
	"net/smtp"
	"sync"
)

type ClientPool interface {
	GetFree() (*smtp.Client, error)
}

type clientPool struct {
	mxs entity.MXs
	wg  *sync.WaitGroup
}

func (p *clientPool) GetFree() (*smtp.Client, error) {
	return &smtp.Client{}, nil
}

func (p *clientPool) lock() {
	p.wg.Add(1)
}

type client struct {
}
