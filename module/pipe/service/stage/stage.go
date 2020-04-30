package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"go.uber.org/fx"
)

const (
	chanSize = 1024
)

type Type int

const (
	UnknownComponentType Type = iota
	SingleComponentType
	MultiComponentType
)

type Constructor func(*entity.Stage, interface{}) (Stage, error)

type Descriptor struct {
	fx.Out
	Name        string
	Type        Type
	Constructor Constructor
}

type Stage interface {
	Run() error
	Bind(Stage)
}

type ResultStage interface {
	Stage
	Results() chan module.Delivery
}

type DeliveryStage interface {
	Stage
	Deliveries() chan module.Delivery
}
