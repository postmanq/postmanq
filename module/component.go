package module

import "go.uber.org/fx"

type ConfigProvider interface {
	Populate(string, interface{}) error
}

type ComponentConstruct func(ConfigProvider) interface{}

type ComponentOut struct {
	fx.Out
	Descriptor ComponentDescriptor `group:"component"`
}

type ComponentDescriptor struct {
	Name      string
	Construct ComponentConstruct
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
	OnReceive(chan Delivery) error
}

type SendComponent interface {
	OnSend(Delivery) error
}

type ProcessComponent interface {
	OnProcess(Delivery) error
}
