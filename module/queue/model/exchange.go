package model

import "github.com/streadway/amqp"

type ExchangeKind string

type Exchange struct {
	Name                       string
	Kind                       string
	Durable                    bool
	AutoDelete                 bool
	Internal                   bool
	NoWait                     bool
	Args                       amqp.Table
	DisablePanicOnChannelError bool
}
