package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"log"
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
					receiver: receiver,
					out:      make(chan module.Delivery, chanSize),
				}, nil
			},
		},
	}
}

type receive struct {
	receiver module.ReceiveComponent
	out      chan module.Delivery
}

func (s *receive) Init() error {
	cmp, ok := s.receiver.(module.InitComponent)
	if ok {
		return cmp.OnInit()
	}

	return nil
}

func (s *receive) Start(in <-chan module.Delivery) <-chan module.Delivery {
	go func() {
		err := s.receiver.OnReceive(s.out)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s.out
}

func (s *receive) Stop() error {
	close(s.out)
	return nil
}
