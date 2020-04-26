package service

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/module/rabbitmq/model"
	"github.com/streadway/amqp"
	"math/rand"
)

type Subscriber interface {
	Subscribe(context.Context) (<-chan amqp.Delivery, error)
	Bind(string, string) error
	Unbind(string, string) error
	Remove() error
	Close() error
}

type subscriber struct {
	channel *amqp.Channel
	queue   model.Queue
	ctx     context.Context
	cancel  context.CancelFunc
	errors  chan *amqp.Error
}

func (s *subscriber) Subscribe(c context.Context) (<-chan amqp.Delivery, error) {
	consumerName := fmt.Sprintf("%s:%d", s.queue.Name, rand.Int())
	return s.channel.Consume(s.queue.Name, consumerName, s.queue.AutoAck, s.queue.Exclusive, false, s.queue.NoWait, s.queue.Args)
}

func (s *subscriber) Bind(key, exchange string) error {
	return s.channel.QueueBind(s.queue.Name, key, exchange, false, nil)
}

func (s *subscriber) Unbind(key, exchange string) error {
	return s.channel.QueueUnbind(s.queue.Name, key, exchange, nil)
}

func (s *subscriber) Remove() error {
	s.cancel()
	_, err := s.channel.QueueDelete(s.queue.Name, false, false, false)
	return err
}

func (s *subscriber) Close() error {
	err := s.channel.Close()
	s.cancel()
	return err
}
