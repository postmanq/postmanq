package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"go.uber.org/multierr"
	"sync"
)

func NewParallelMiddleware() Descriptor {
	return Descriptor{
		Name: "parallel_middleware",
		Type: MultiComponentType,
		Constructor: func(cfg *entity.Stage, rawComponents interface{}) (Stage, error) {
			components := rawComponents.([]interface{})
			middlewares := make([]module.ProcessComponent, len(components))
			for i, component := range components {
				middleware, ok := component.(module.ProcessComponent)
				if !ok {
					return nil, errors.CantCastTypeToComponent(middleware)
				}

				middlewares[i] = middleware
			}

			return &parallelMiddleware{
				deliveries:  make(chan module.Delivery, chanSize),
				middlewares: middlewares,
			}, nil
		},
	}
}

type parallelMiddleware struct {
	deliveries  chan module.Delivery
	middlewares []module.ProcessComponent
	prev        ResultStage
	next        DeliveryStage
}

func (s *parallelMiddleware) Run() error {
	defer func() {
		close(s.deliveries)
	}()

	numSenders := len(s.middlewares)
	for delivery := range s.deliveries {
		wg := new(sync.WaitGroup)
		wg.Add(numSenders)

		for y := 0; y < numSenders; y++ {
			go func(middleware module.ProcessComponent) {
				err := middleware.OnProcess(delivery)
				if err != nil {
					delivery.Err = multierr.Append(delivery.Err, err)
				}
				wg.Done()
			}(s.middlewares[y])
		}

		wg.Wait()

		if delivery.Err == nil {
			s.next.Deliveries() <- delivery
		} else {
			s.prev.Results() <- delivery
		}
	}

	return nil
}

func (s *parallelMiddleware) Deliveries() chan module.Delivery {
	return s.deliveries
}

func (s *parallelMiddleware) Bind(any Stage) {
	switch a := any.(type) {
	case ResultStage:
		s.prev = a
	case DeliveryStage:
		s.next = a
	}
}
