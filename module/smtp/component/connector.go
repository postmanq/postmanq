package component

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/smtp/service"
)

type connector struct {
	scanner service.Scanner
}

func (c *connector) OnProcess(delivery module.Delivery) error {

	return nil
}
