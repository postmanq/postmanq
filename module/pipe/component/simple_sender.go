package component

import (
	"github.com/byorty/postmanq_sender/plugin"
	"github.com/byorty/postmanq_sender/sending"
)

type SimpleSender struct {
	sendings chan sending.Sending
	sender   plugin.Sender
	next     SendingStep
}

func NewSimpleSender(sender plugin.Sender, next SendingStep) *SimpleSender {
	return &SimpleSender{
		sendings: make(chan sending.Sending, 1024),
		sender:   sender,
		next:     next,
	}
}

func (c *SimpleSender) Run(numCPU int) {
	for i := 0; i < numCPU; i++ {
		go func() {
			for s := range c.sendings {
				err := c.sender.OnSend(s)

				if err == nil {
					c.next.Send(s)
				} else {
					s.Abort(err)
				}
			}
		}()
	}
}

func (c *SimpleSender) Send(s sending.Sending) {
	c.sendings <- s
}
