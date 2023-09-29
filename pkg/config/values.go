package config

import "go.uber.org/zap"

type Config struct {
	Logger zap.Config `yaml:"log"`
}
