package component

import (
	"github.com/postmanq/postmanq/module"
	"sync"
)

type ParallelSender struct {
	deliveries chan module.Delivery
	senders    []module.SendComponent
	prev       ResultStep
	next       DeliveryStep
}

func NewParallelSender(senders []module.SendComponent, prev ResultStep, next DeliveryStep) *ParallelSender {
	return &ParallelSender{
		deliveries: make(chan module.Delivery, 1024),
		senders:    senders,
		prev:       prev,
		next:       next,
	}
}

func (c *ParallelSender) Run(workerCount int) {
	numSenders := len(c.senders)

	for i := 0; i < workerCount; i++ {
		go func() {
			//senderErrors := concurrent.NewSlice()
			for s := range c.deliveries {
				wg := new(sync.WaitGroup)
				wg.Add(numSenders)

				for y := 0; y < numSenders; y++ {
					go func(sender module.SendComponent, delivery module.Delivery) {
						err := sender.OnSend(delivery)
						if err != nil {
							//senderErrors.Append(err)
						}
						wg.Done()
					}(c.senders[y], s)
				}

				wg.Wait()

				//if senderErrors.Len() == 0 {
				//	c.next.Send(s)
				//} else {
				//	var err error
				//	for i := 0; i < senderErrors.Len(); i++ {
				//		err = multierr.Append(err, senderErrors.Get(i).(error))
				//	}
				//
				//	s.Abort(err)
				//}
				//
				//senderErrors.Clear()
			}
		}()
	}
}
