package stage

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/pipe/entity"
	"go.uber.org/fx"
)

const (
	chanSize = 1024
)

type ArgType int

const (
	ArgTypeUnknown ArgType = iota
	ArgTypeSingle
	ArgTypeMulti
)

type Constructor func(*entity.Stage, interface{}) (Stage, error)

type Out struct {
	fx.Out
	Descriptor `group:"stage"`
}

type Descriptor struct {
	Name        string
	Type        ArgType
	Constructor Constructor
}

type Stage interface {
	Init() error
	Start(<-chan module.Delivery) <-chan module.Delivery
	Stop() error
}
