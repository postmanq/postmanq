package model

import "github.com/postmanq/postmanq/module"

type Component struct {
	Name   string                 `yaml:"name"`
	Config module.ComponentConfig `yaml:"config"`
}
