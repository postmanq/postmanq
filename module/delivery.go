package module

import "net/smtp"

type Email struct {
	Sender        string
	Recipient     string
	RecipientHost string
	Client        *smtp.Client
}

type Delivery struct {
	Email *Email
	Err   chan error
}

func (d Delivery) Cancel(err error) {
	d.Err <- err
}

func (d Delivery) Complete() {
	d.Err <- nil
}

func (d Delivery) Next(next chan<- Delivery) {
	next <- d
}
