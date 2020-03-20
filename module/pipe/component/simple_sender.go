package component

import (
	"github.com/postmanq/postmanq/module"
)

type SimpleSender struct {
	deliveries chan module.Delivery
	sender     module.SendComponent
	prev       ResultStep
	next       DeliveryStep
}

func NewSimpleSender(sender module.SendComponent, prev ResultStep, next DeliveryStep) *SimpleSender {
	return &SimpleSender{
		deliveries: make(chan module.Delivery, 1024),
		sender:     sender,
		prev:       prev,
		next:       next,
	}
}

func (c *SimpleSender) Run(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for delivery := range c.deliveries {
				err := c.sender.OnSend(delivery)

				if err == nil {
					c.next.Deliveries() <- delivery
				} else {
					c.prev.Results() <- delivery
				}
			}
		}()
	}
}
