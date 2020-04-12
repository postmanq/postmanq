package module

type Component interface {
	GetName() string
}

type InitComponent interface {
	Component
	OnInit() error
}

type ReceiveComponent interface {
	Component
	OnReceive(chan Delivery, chan Delivery) error
}

type SendComponent interface {
	Component
	OnSend(Delivery) error
}

type ProcessComponent interface {
	Component
	OnProcess(Delivery) error
}
