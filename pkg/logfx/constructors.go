package logfx

import (
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/logfx/internal/services"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"go.uber.org/fx"
	"runtime"
)

var (
	Constructors = fx.Provide(
		NewFxAppLogger,
	)
)

func NewFxAppLogger(cfg *config.Config) (log.Logger, error) {
	logger, err := services.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}
	logger.Debugf("app start cpu: %d", runtime.NumCPU())
	return logger, nil
}
