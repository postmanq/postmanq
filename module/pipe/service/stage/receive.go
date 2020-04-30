package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
)

func NewReceive() Out {
	return Out{
		Descriptor: Descriptor{
			Name: "receive",
			Type: ArgTypeSingle,
			Constructor: func(cfg *entity.Stage, component interface{}) (Stage, error) {
				receiver, ok := component.(module.ReceiveComponent)
				if !ok {
					return nil, errors.CantCastTypeToComponent(component)
				}

				return &receive{
					receiver:   receiver,
					deliveries: make(chan module.Delivery, chanSize),
					results:    make(chan module.Delivery, chanSize),
				}, nil
			},
		},
	}
}

type receive struct {
	receiver   module.ReceiveComponent
	next       DeliveryStage
	deliveries chan module.Delivery
	results    chan module.Delivery
}

func (s *receive) Start() error {
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

func (s *receive) Results() chan module.Delivery {
	return s.results
}

func (s *receive) Bind(next Stage) {
	s.next = next.(DeliveryStage)
}
