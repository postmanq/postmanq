package module

type Delivery struct {
	Err chan error
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
