package component

import (
	"github.com/postmanq/postmanq/module"
)

type ParallelReceiver struct {
	receivers []module.ReceiveComponent
	next      DeliveryStep
}

func NewParallelReceiver(receivers []module.ReceiveComponent, next DeliveryStep) *ParallelReceiver {
	return &ParallelReceiver{
		receivers: receivers,
		next:      next,
	}
}

func (r *ParallelReceiver) Run(workerCount int) {
	for _, receiver := range r.receivers {
		deliveries := make(chan module.Delivery, workerCount)
		results := make(chan module.Delivery, workerCount)
		go func() {
			err := receiver.OnReceive(deliveries, results)
			if err != nil {
				close(deliveries)
				close(results)
			}
		}()

		for i := 0; i < workerCount; i++ {
			go func() {
				for delivery := range deliveries {
					r.next.Deliveries() <- delivery
				}

				close(deliveries)
				close(results)
			}()
		}
	}
}
