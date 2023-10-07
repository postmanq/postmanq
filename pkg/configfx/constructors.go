package configfx

import (
	config2 "github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/configfx/internal/services"
	"go.uber.org/fx"
)

var (
	Constructors = fx.Provide(
		services.NewFxConfigProviderFactory,
		NewFxConfigProvider,
		NewFxConfig,
	)
)

func NewFxConfigProvider(
	factory config2.ProviderFactory,
	filename string,
) (config2.Provider, error) {
	return factory.Create(config2.File(filename))
}

func NewFxConfig(
	provider config2.Provider,
) (*config2.Config, error) {
	cfg := new(config2.Config)
	err := provider.Populate(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
