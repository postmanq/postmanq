package postmanqfx

import (
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"postmanq",
		fx.Provide(
			NewFxConfig,
		),
	)
)

func NewFxConfig(
	provider config.Provider,
) (*postmanq.Config, error) {
	cfg := new(postmanq.Config)
	err := provider.Populate(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
