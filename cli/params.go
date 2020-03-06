package cli

import (
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Components []interface{} `group:"component"`
}
