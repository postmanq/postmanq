package component

import "github.com/postmanq/postmanq/module"

type Step interface {
	Run(int)
}

type ResultStep interface {
	Results() chan<- module.Delivery
}

type DeliveryStep interface {
	Step
	Deliveries() chan<- module.Delivery
}
