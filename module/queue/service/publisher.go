package service

import (
	"context"
	"github.com/postmanq/postmanq/module/queue/model"
	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(string, amqp.Publishing) error
	Bind(string, string) error
	Unbind(string, string) error
	Close() error
}

type publisher struct {
	exchange  model.Exchange
	channel   *amqp.Channel
	cancel    context.CancelFunc
	ctx       context.Context
	errorChan chan *amqp.Error
}

func (p *publisher) Publish(key string, msg amqp.Publishing) error {
	return p.channel.Publish(p.exchange.Name, key, false, false, msg)
}

func (p *publisher) Bind(key string, exchange string) error {
	return p.channel.ExchangeBind(p.exchange.Name, key, exchange, false, nil)
}

func (p *publisher) Unbind(key string, exchange string) error {
	return p.channel.ExchangeUnbind(p.exchange.Name, key, exchange, false, nil)
}

func (p *publisher) Close() error {
	err := p.channel.Close()
	p.cancel()
	return err
}
