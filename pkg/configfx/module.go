package configfx

import (
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/configfx/internal/services"
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
	filename string,
) (config.Provider, error) {
	return factory.Create(config.File(filename))
}
