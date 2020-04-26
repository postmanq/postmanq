package module

import "go.uber.org/fx"

type ComponentConfig map[string]interface{}

type ComponentConstruct func(ComponentConfig) interface{}

type ComponentDescriptor struct {
	fx.Out
	Name      string
	Construct ComponentConstruct `group:"component"`
}

type ComponentType int

const (
	ComponentTypeInit ComponentType = iota
	ComponentTypeReceive
	ComponentTypeSend
	ComponentTypeProcess
)

type InitComponent interface {
	OnInit() error
}

type ReceiveComponent interface {
	OnReceive(chan Delivery, chan Delivery) error
}

type SendComponent interface {
	OnSend(Delivery) error
}

type ProcessComponent interface {
	OnProcess(Delivery) error
}
