package service

import "github.com/postmanq/postmanq/module"

type ReceiveStage struct {
	receiver   module.ReceiveComponent
	next       DeliveryStage
	deliveries chan module.Delivery
	results    chan module.Delivery
}

func NewReceiveStage(receiver module.ReceiveComponent) *ReceiveStage {
	return &ReceiveStage{
		receiver:   receiver,
		deliveries: make(chan module.Delivery, chanSize),
		results:    make(chan module.Delivery, chanSize),
	}
}

func (s *ReceiveStage) Run() error {
	defer func() {
		close(s.deliveries)
		close(s.results)
	}()

	go func() {
		for delivery := range s.deliveries {
			s.next.Deliveries() <- delivery
		}
	}()

	return s.receiver.OnReceive(s.deliveries, s.results)
}

func (s *ReceiveStage) Results() chan module.Delivery {
	return s.results
}

func (s *ReceiveStage) Bind(next DeliveryStage) {
	s.next = next
}
