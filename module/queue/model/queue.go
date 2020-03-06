package model

import "github.com/streadway/amqp"

type Queue struct {
	Name                       string
	Exchange                   string
	Key                        string
	Keys                       []string
	Durable                    bool
	AutoDelete                 bool
	Exclusive                  bool
	NoWait                     bool
	PrefetchCount              int
	Args                       amqp.Table
	AutoAck                    bool
	ExchangeKeyMap             map[string][]string
	DisablePanicOnChannelError bool
}
