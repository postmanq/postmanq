package logfx

import (
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/internal/services"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	cfg := new(zap.Config)
	err := provider.PopulateByKey("log", cfg)
	if err != nil {
		return nil, err
	}

	logger, err := services.NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	logger.Debugf("app start cpu: %d", runtime.NumCPU())
	return logger, nil
}
