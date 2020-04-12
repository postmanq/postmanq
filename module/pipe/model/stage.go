package model

type StageType string

const (
	StageTypeReceive            StageType = "receive"
	StageTypeMiddleware         StageType = "middleware"
	StageTypeParallelMiddleware StageType = "parallel_middleware"
	StageTypeComplete           StageType = "complete"
)

type Stage struct {
	Type       StageType `yaml:"type"`
	Component  string    `yaml:"component"`
	Components []string  `yaml:"components"`
}
