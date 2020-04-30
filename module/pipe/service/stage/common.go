package stage

import "github.com/postmanq/postmanq/module"

type common struct {
	deliveries chan module.Delivery
	cmp        interface{}
}

func (s *common) Init() error {
	cmp, ok := s.cmp.(module.InitComponent)
	if ok {
		return cmp.OnInit()
	}

	return nil
}

func (s *common) Start() error {
	return nil
}

func (s *common) Stop() error {
	close(s.deliveries)
	return nil
}

func (s *common) Deliveries() chan module.Delivery {
	return s.deliveries
}
