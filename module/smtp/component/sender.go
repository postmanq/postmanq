package component

import (
	"fmt"
	"github.com/postmanq/postmanq/module"
	"github.com/postmanq/postmanq/module/smtp/service"
)

type sender struct {
	scanner        service.Scanner
	connector      service.Connector
	configProvider module.ConfigProvider
}

func NewSender() module.ComponentOut {
	return module.ComponentOut{
		Descriptor: module.ComponentDescriptor{
			Name: "smtp/sender",
			Construct: func(configProvider module.ConfigProvider) interface{} {
				return &sender{
					configProvider: configProvider,
				}
			},
		},
	}
}

func (c *sender) OnSend(delivery module.Delivery) error {
	result := c.scanner.Scan(delivery.Email.RecipientHost)
	if result.GetStatus() != service.ScannerResultStatusSuccess {
		return result.GetError()
	}

	pool := c.connector.Connect(result)
	client, err := pool.GetFree()
	if err != nil {
		return err
	}

	//client.SetTimeout(common.App.Timeout().Mail)
	err = client.Mail(delivery.Email.Sender)
	if err != nil {
		return err
	}

	//client.SetTimeout(common.App.Timeout().Rcpt)
	err = client.Rcpt(delivery.Email.Recipient)
	if err != nil {
		return err
	}

	//client.SetTimeout(common.App.Timeout().Data)
	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(wc, delivery.Email.Body)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	err = client.Reset()
	if err != nil {
		return err
	}

	return nil
}
