package cli

import (
	"github.com/postmanq/postmanq/module/pipe/component"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Pipe *component.Runner
}
