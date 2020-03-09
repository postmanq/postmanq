package module

type InitComponent interface {
	OnInit() error
}

type ReceiveComponent interface {
	OnReceive(chan Delivery, chan Delivery) error
}

type SendComponent interface {
	OnSend() error
}
