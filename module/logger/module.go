package main

import "go.uber.org/zap"

func PqModule() *zap.SugaredLogger {
	logger := zap.NewExample().Sugar()
	defer logger.Sync()
	return logger
}
