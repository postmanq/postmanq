package component

import "github.com/postmanq/postmanq/module"

type Step interface {
	Run(int)
}

type SendingStep interface {
	Step
	Deliveries() chan<- module.Delivery
	Results() chan<- module.Delivery
}
