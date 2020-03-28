package model

type StageType string

type Stage struct {
	Type       StageType `yaml:"type"`
	Component  string    `yaml:"component"`
	Components []string  `yaml:"components"`
}
