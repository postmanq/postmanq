package component

import (
	"github.com/postmanq/postmanq/module"
)

type CompleteSender struct {
	deliveries chan module.Delivery
	sender     module.SendComponent
	prev       ResultStep
}

func NewCompleteSender(sender module.SendComponent, prev ResultStep) *CompleteSender {
	return &CompleteSender{
		deliveries: make(chan module.Delivery, 1024),
		sender:     sender,
		prev:       prev,
	}
}

func (c *CompleteSender) Run(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for delivery := range c.deliveries {
				err := c.sender.OnSend(delivery)

				if err == nil {
					//delivery.Complete()
				} else {
					//delivery.Abort(err)
				}
			}
		}()
	}
}
