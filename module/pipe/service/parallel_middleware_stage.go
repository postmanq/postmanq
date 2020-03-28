package service

import (
	"github.com/postmanq/postmanq/module"
	"go.uber.org/multierr"
	"sync"
)

type ParallelMiddlewareStage struct {
	deliveries  chan module.Delivery
	middlewares []module.ProcessComponent
	prev        ResultStage
	next        DeliveryStage
}

func NewParallelMiddlewareStage(middlewares []module.ProcessComponent, prev ResultStage) *ParallelMiddlewareStage {
	return &ParallelMiddlewareStage{
		deliveries:  make(chan module.Delivery, chanSize),
		middlewares: middlewares,
		prev:        prev,
	}
}

func (c *ParallelMiddlewareStage) Run() error {
	defer func() {
		close(c.deliveries)
	}()

	numSenders := len(c.middlewares)
	for delivery := range c.deliveries {
		wg := new(sync.WaitGroup)
		wg.Add(numSenders)

		for y := 0; y < numSenders; y++ {
			go func(middleware module.ProcessComponent) {
				err := middleware.OnProcess(delivery)
				if err != nil {
					delivery.Err = multierr.Append(delivery.Err, err)
				}
				wg.Done()
			}(c.middlewares[y])
		}

		wg.Wait()

		if delivery.Err == nil {
			c.next.Deliveries() <- delivery
		} else {
			c.prev.Results() <- delivery
		}
	}

	return nil
}

func (c *ParallelMiddlewareStage) Deliveries() chan module.Delivery {
	return c.deliveries
}

func (c *ParallelMiddlewareStage) Bind(next DeliveryStage) {
	c.next = next
}
