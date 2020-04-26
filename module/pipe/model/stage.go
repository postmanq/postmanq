package model

type StageType string

const (
	StageTypeReceive            StageType = "receive"
	StageTypeMiddleware         StageType = "middleware"
	StageTypeParallelMiddleware StageType = "parallel_middleware"
	StageTypeComplete           StageType = "complete"
)

type Stage struct {
	Name       string      `yaml:"name"`
	Type       StageType   `yaml:"type"`
	Component  Component   `yaml:"component"`
	Components []Component `yaml:"components"`
}
