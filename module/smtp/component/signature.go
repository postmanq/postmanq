package component

import "github.com/postmanq/postmanq/module"

type signature struct {
}

func (c *signature) OnProcess(delivery module.Delivery) error {
	return nil
}
