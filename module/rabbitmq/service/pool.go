package service

import (
	"context"
	"github.com/postmanq/postmanq/module/rabbitmq/model"
	"github.com/streadway/amqp"
	"sync"
)

type Pool interface {
	Connect([]string) error
	CreatePublisher(context.Context, model.Exchange) (Publisher, error)
	CreateSubscriber(context.Context, model.Queue) (Subscriber, error)
}

type pool struct {
	wMx        *sync.RWMutex
	wSelector  int
	writeConns []*amqp.Connection
	rMx        *sync.RWMutex
	rSelector  int
	readConns  []*amqp.Connection
}

func NewPool() Pool {
	p := &pool{
		readConns:  make([]*amqp.Connection, 0),
		writeConns: make([]*amqp.Connection, 0),
		rSelector:  0,
		wSelector:  0,
		rMx:        new(sync.RWMutex),
		wMx:        new(sync.RWMutex),
	}
	return p
}

func (p *pool) Connect(urls []string) error {
	n := len(urls)
	connPool := make([]*amqp.Connection, n)

	for i, url := range urls {
		conn, err := amqp.Dial(url)
		if err != nil {
			return err
		}

		connPool[i] = conn
	}

	if n >= 2 {
		half := n / 2
		p.writeConns = connPool[:half]
		p.readConns = connPool[half:]
	} else {
		p.writeConns = connPool
		p.readConns = connPool
	}

	return nil
}

func (p *pool) readConn() *amqp.Connection {
	if p.readConns != nil {
		p.rMx.Lock()
		q := p.readConns[p.rSelector]

		if p.rSelector == len(p.readConns)-1 {
			p.rSelector = 0
		} else {
			p.rSelector++
		}

		p.rMx.Unlock()
		return q
	}
	return nil
}

func (p *pool) writeConn() *amqp.Connection {
	if p.writeConns != nil {
		p.wMx.Lock()
		q := p.writeConns[p.rSelector]

		if p.wSelector == len(p.writeConns)-1 {
			p.wSelector = 0
		} else {
			p.wSelector++
		}

		p.wMx.Unlock()
		return q
	}
	return nil
}

func (p *pool) CreatePublisher(ctx context.Context, exchange model.Exchange) (Publisher, error) {
	channel, err := p.writeConn().Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(exchange.Name, exchange.Kind, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, exchange.Args)
	if err != nil {
		return nil, err
	}

	chCtx, cancel := context.WithCancel(ctx)
	pub := &publisher{
		exchange:  exchange,
		channel:   channel,
		cancel:    cancel,
		ctx:       chCtx,
		errorChan: make(chan *amqp.Error),
	}
	if !exchange.DisablePanicOnChannelError {
		go p.listenErrors(chCtx, channel, pub.errorChan)
	}
	return pub, nil
}

func (p *pool) CreateSubscriber(ctx context.Context, q model.Queue) (Subscriber, error) {
	channel, err := p.readConn().Channel()
	if err != nil {
		return nil, err
	}

	err = channel.Qos(q.PrefetchCount, 0, false)
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.Args)
	if err != nil {
		return nil, err
	}

	chCtx, cancel := context.WithCancel(ctx)
	s := &subscriber{
		channel: channel,
		queue:   q,
		ctx:     chCtx,
		cancel:  cancel,
		errors:  make(chan *amqp.Error),
	}

	if len(q.ExchangeKeyMap) > 0 {
		for ex, keyMap := range q.ExchangeKeyMap {
			for _, key := range keyMap {
				err = s.Bind(key, ex)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		if len(q.Keys) > 0 {
			for _, key := range q.Keys {
				err = s.Bind(key, q.Exchange)
				if err != nil {
					return nil, err
				}
			}
		} else {
			err = s.Bind(q.Key, q.Exchange)
			if err != nil {
				return nil, err
			}
		}
	}

	if !s.queue.DisablePanicOnChannelError {
		go p.listenErrors(chCtx, s.channel, s.errors)
	}

	return s, nil
}

func (p *pool) listenErrors(ctx context.Context, channel *amqp.Channel, errorChan chan *amqp.Error) {
	channel.NotifyClose(errorChan)
	for err := range errorChan {
		select {
		case <-ctx.Done():
			return
		default:
			if err != nil {
				panic(err)
			}
		}
	}
}
