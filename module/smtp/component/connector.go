package component

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/smtp/service/scanner"
)

type connector struct {
	scanner scanner.Scanner
}

func (c *connector) OnProcess(delivery module.Delivery) error {
	result := c.scanner.Scan(delivery.Email.RecipientHost)
	if result.GetStatus() != scanner.ResultStatusSuccess {
		return result.GetError()
	}

	return nil
}
