package component

import "github.com/postmanq/postmanq/module"

type MiddlewareStage struct {
	deliveries chan module.Delivery
	middleware module.ProcessComponent
	prev       ResultStage
	next       DeliveryStage
}

func NewMiddlewareStage(middleware module.ProcessComponent, prev ResultStage) *MiddlewareStage {
	return &MiddlewareStage{
		deliveries: make(chan module.Delivery, chanSize),
		middleware: middleware,
		prev:       prev,
	}
}

func (s *MiddlewareStage) Run() error {
	defer func() {
		close(s.deliveries)
	}()

	for delivery := range s.deliveries {
		delivery.Err = s.middleware.OnProcess(delivery)

		if delivery.Err == nil {
			s.next.Deliveries() <- delivery
		} else {
			s.prev.Results() <- delivery
		}
	}

	return nil
}

func (s *MiddlewareStage) Deliveries() chan module.Delivery {
	return s.deliveries
}

func (s *MiddlewareStage) Bind(next DeliveryStage) {
	s.next = next
}
