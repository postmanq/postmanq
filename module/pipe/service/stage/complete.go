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
					sender: sender,
				}, nil
			},
		},
	}
}

type complete struct {
	sender module.SendComponent
}

func (s *complete) Init() error {
	cmp, ok := s.sender.(module.InitComponent)
	if ok {
		return cmp.OnInit()
	}

	return nil
}

func (s *complete) Start(in <-chan module.Delivery) <-chan module.Delivery {
	go func() {
		for delivery := range in {
			err := s.sender.OnSend(delivery)
			if err == nil {
				delivery.Complete()
			} else {
				delivery.Cancel(err)
			}
		}
	}()

	return nil
}

func (s *complete) Stop() error {
	return nil
}
