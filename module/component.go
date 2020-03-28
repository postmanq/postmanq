package module

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
