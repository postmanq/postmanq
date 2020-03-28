package component

import "github.com/postmanq/postmanq/module"

const (
	chanSize = 1024
)

type Stage interface {
	Run() error
	Bind(DeliveryStage)
}

type ResultStage interface {
	Stage
	Results() chan module.Delivery
}

type DeliveryStage interface {
	Stage
	Deliveries() chan module.Delivery
}
