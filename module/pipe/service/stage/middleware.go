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
					out:        make(chan module.Delivery, chanSize),
					middleware: m,
				}, nil
			},
		},
	}
}

type middleware struct {
	out        chan module.Delivery
	middleware module.ProcessComponent
}

func (s *middleware) Init() error {
	cmp, ok := s.middleware.(module.InitComponent)
	if ok {
		return cmp.OnInit()
	}

	return nil
}

func (s *middleware) Start(in <-chan module.Delivery) <-chan module.Delivery {
	go func() {
		for delivery := range in {
			err := s.middleware.OnProcess(delivery)
			if err == nil {
				delivery.Next(s.out)
			} else {
				delivery.Cancel(err)
			}
		}
	}()

	return s.out
}

func (s *middleware) Stop() error {
	close(s.out)
	return nil
}
