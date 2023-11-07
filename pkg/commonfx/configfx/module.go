package configfx

import (
	"github.com/postmanq/postmanq/pkg/commonfx/app"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/internal/services"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"config",
		fx.Provide(
			services.NewFxConfigProviderFactory,
			NewFxConfigProvider,
		),
	)
)

func NewFxConfigProvider(
	factory config.ProviderFactory,
	args app.Arguments,
) (config.Provider, error) {
	return factory.Create(config.File(args.ConfigFilename))
}
