package main

import (
	"github.com/postmanq/postmanq/pkg/config"
	"github.com/postmanq/postmanq/pkg/log"
	"github.com/postmanq/postmanq/pkg/log/internal/services"
	"github.com/postmanq/postmanq/pkg/plugin"
	"runtime"
)

var (
	Plugin plugin.Construct = func() plugin.Plugin {
		return plugin.Plugin{
			Constructs: []interface{}{
				NewFxAppLogger,
			},
		}
	}
)

func NewFxAppLogger(cfg *config.Config) (log.Logger, error) {
	logger, err := services.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}
	logger.Debugf("app start cpu: %d", runtime.NumCPU())
	return logger, nil
}
