package service

import (
	"context"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(interface{}) error
	StructCtx(context.Context, interface{}) error
	Var(interface{}, string) error
	VarCtx(context.Context, interface{}, string) error
}

func NewValidator() Validator {
	return validator.New()
}
