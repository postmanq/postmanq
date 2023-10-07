package config

import "go.uber.org/zap"

type Config struct {
	NumCPU  int        `yaml:"num_cpu"`
	Logger  zap.Config `yaml:"log"`
	Plugins []Plugin   `yaml:"plugins"`
}

type Plugin struct {
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config"`
	Tags   []string    `yaml:"tags"`
}
