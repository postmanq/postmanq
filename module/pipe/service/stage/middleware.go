package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
)

func NewMiddleware() Out {
	return Out{
		Descriptor: Descriptor{
			Name: "middleware",
			Type: ArgTypeSingle,
			Constructor: func(cfg *entity.Stage, component interface{}) (Stage, error) {
				m, ok := component.(module.ProcessComponent)
				if !ok {
					return nil, errors.CantCastTypeToComponent(component)
				}

				return &middleware{
					deliveries: make(chan module.Delivery, chanSize),
					middleware: m,
				}, nil
			},
		},
	}
}

type middleware struct {
	deliveries chan module.Delivery
	middleware module.ProcessComponent
	prev       ResultStage
	next       DeliveryStage
}

func (s *middleware) Start() error {
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

func (s *middleware) Deliveries() chan module.Delivery {
	return s.deliveries
}

func (s *middleware) Bind(any Stage) {
	switch a := any.(type) {
	case ResultStage:
		s.prev = a
	case DeliveryStage:
		s.next = a
	}
}
