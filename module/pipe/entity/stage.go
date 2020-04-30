package entity

type Stage struct {
	Name       string       `yaml:"name"`
	Type       string       `yaml:"type"`
	Component  *Component   `yaml:"component"`
	Components []*Component `yaml:"components"`
}
