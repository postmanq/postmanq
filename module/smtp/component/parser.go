package component

import "github.com/postmanq/postmanq/module"

type parser struct {
}

func (c *parser) OnProcess(delivery module.Delivery) error {
	return nil
}
