package main

import (
	"github.com/postmanq/postmanq/cli"
	"github.com/postmanq/postmanq/module/config/service"
	"go.uber.org/fx"
	"log"
)

type PqModuleOut struct {
	fx.Out
	ConfigProvider service.ConfigProvider
}

func PqModule(args cli.Arguments) (PqModuleOut, error) {
	configProvider, err := service.NewConfigProviderByFile(args.ConfigFilename)

	var out PqModuleOut
	out.ConfigProvider = configProvider

	log.Println("call config module")

	return out, err
}
