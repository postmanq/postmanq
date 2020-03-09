package module

type InitComponent interface {
	OnInit() error
}

type ReceiveComponent interface {
	OnReceive() error
}
