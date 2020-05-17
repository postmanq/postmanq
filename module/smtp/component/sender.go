package component

import (
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/smtp/service"
)

type sender struct {
	scanner service.Scanner
}

func (c *sender) OnSend(delivery module.Delivery) error {
	result := c.scanner.Scan(delivery.Email.RecipientHost)
	if result.GetStatus() != service.ScannerResultStatusSuccess {
		return result.GetError()
	}

	return nil
}
