package component

import (
	"github.com/postmanq/postmanq/module"
)

type ParallelReceiver struct {
	receivers []module.ReceiveComponent
	next      SendingStep
}

func NewParallelReceiver(receivers []module.ReceiveComponent, next SendingStep) *ParallelReceiver {
	return &ParallelReceiver{
		receivers: receivers,
		next:      next,
	}
}

func (r *ParallelReceiver) Run(numCPU int) {
	for _, receiver := range r.receivers {
		for i := 0; i < numCPU; i++ {
			go func() {
				deliveries := make(chan module.Delivery, numCPU)
				results := make(chan module.Delivery, numCPU)

				go func() {
					err := receiver.OnReceive(deliveries, results)
					if err != nil {
						close(deliveries)
						close(results)
					}
				}()
				for delivery := range deliveries {
					r.next.Deliveries() <- delivery
				}

				close(deliveries)
				close(results)
			}()
		}
	}
}
