package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/postmanq/postmanq/module/validator/service"
	"go.uber.org/fx"
)

type PqModuleOut struct {
	fx.Out
	Validator service.Validator
}

func PqModule() PqModuleOut {
	return PqModuleOut{
		Validator: validator.New(),
	}
}
