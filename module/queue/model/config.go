package model

import "time"

type Config struct {
	Url                  string          `yaml:"url" validate:"required"`
	PrefetchMessageCount int             `yaml:"prefetchMessageCount" validate:"required"`
	Repeats              []time.Duration `yaml:"repeats" validate:"required"`
}
