package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"github.com/postmanq/postmanq/module/pipe/errors"
	"go.uber.org/multierr"
	"sync"
)

func NewParallelMiddleware() Out {
	return Out{
		Descriptor: Descriptor{
			Name: "parallel_middleware",
			Type: ArgTypeMulti,
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
					out:         make(chan module.Delivery, chanSize),
					middlewares: middlewares,
				}, nil
			},
		},
	}
}

type parallelMiddleware struct {
	out         chan module.Delivery
	middlewares []module.ProcessComponent
}

func (s *parallelMiddleware) Init() error {
	for _, middleware := range s.middlewares {
		cmp, ok := middleware.(module.InitComponent)
		if ok {
			return cmp.OnInit()
		}
	}

	return nil
}

func (s *parallelMiddleware) Start(in <-chan module.Delivery) <-chan module.Delivery {
	middlewareLen := len(s.middlewares)
	go func() {
		for delivery := range in {
			wg := new(sync.WaitGroup)
			wg.Add(middlewareLen)

			var multiErr error
			for y := 0; y < middlewareLen; y++ {
				go func(middleware module.ProcessComponent) {
					err := middleware.OnProcess(delivery)
					if err != nil {
						multiErr = multierr.Append(multiErr, err)
					}
					wg.Done()
				}(s.middlewares[y])
			}

			wg.Wait()
			if multiErr == nil {
				delivery.Next(s.out)
			} else {
				delivery.Cancel(multiErr)
			}
		}
	}()

	return s.out
}

func (s *parallelMiddleware) Stop() error {
	close(s.out)
	return nil
}
