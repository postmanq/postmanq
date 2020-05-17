package entity

type Component struct {
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config"`
}
