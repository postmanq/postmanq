package component

import (
	"github.com/postmanq/postmanq/module"
)

type CompleteStage struct {
	deliveries chan module.Delivery
	sender     module.SendComponent
	prev       ResultStage
}

func NewCompleteStage(sender module.SendComponent, prev ResultStage) *CompleteStage {
	return &CompleteStage{
		deliveries: make(chan module.Delivery, chanSize),
		sender:     sender,
		prev:       prev,
	}
}

func (s *CompleteStage) Run() error {
	defer func() {
		close(s.deliveries)
	}()

	for delivery := range s.deliveries {
		delivery.Err = s.sender.OnSend(delivery)
		s.prev.Results() <- delivery
	}

	return nil
}

func (s *CompleteStage) Deliveries() chan module.Delivery {
	return s.deliveries
}

func (s *CompleteStage) Bind(next DeliveryStage) {}
