package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
)

func NewComplete() Out {
	return Out{
		Descriptor: Descriptor{
			Name: "complete",
			Type: ArgTypeSingle,
			Constructor: func(cfg *entity.Stage, component interface{}) (Stage, error) {
				sender, ok := component.(module.SendComponent)
				if !ok {
					return nil, errors.CantCastTypeToComponent(component)
				}

				return &complete{
					deliveries: make(chan module.Delivery, chanSize),
					sender:     sender,
				}, nil
			},
		},
	}
}

type complete struct {
	deliveries chan module.Delivery
	sender     module.SendComponent
	prev       ResultStage
}

func (s *complete) Start() error {
	for delivery := range s.deliveries {
		delivery.Err = s.sender.OnSend(delivery)
		s.prev.Results() <- delivery
	}

	return nil
}

func (s *complete) Stop() {
	close(s.deliveries)
}

func (s *complete) Deliveries() chan module.Delivery {
	return s.deliveries
}

func (s *complete) Bind(prev Stage) {
	s.prev = prev.(ResultStage)
}
