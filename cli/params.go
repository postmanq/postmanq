package cli

import (
	"github.com/postmanq/postmanq/module/config/service"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	ConfigProvider service.ConfigProvider
}
