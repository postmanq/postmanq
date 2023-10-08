package logfx

import (
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/logfx/internal/services"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"go.uber.org/fx"
	"runtime"
)

var (
	Module = fx.Module(
		"log",
		fx.Provide(
			NewFxAppLogger,
		),
	)
)

func NewFxAppLogger(provider config.Provider) (log.Logger, error) {
	cfg := new(log.Config)
	err := provider.Populate(cfg)
	if err != nil {
		return nil, err
	}

	logger, err := services.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	logger.Debugf("app start cpu: %d", runtime.NumCPU())
	return logger, nil
}
