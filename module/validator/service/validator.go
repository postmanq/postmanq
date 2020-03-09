package service

import "context"

type Validator interface {
	Struct(interface{}) error
	StructCtx(context.Context, interface{}) error
	Var(interface{}, string) error
	VarCtx(context.Context, interface{}, string) error
}
