package model

type Pipe struct {
	Name     string  `yaml:"name"`
	Replicas int     `yaml:"replicas"`
	Stages   []Stage `yaml:"stages"`
}
