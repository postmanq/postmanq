package template

import "github.com/pkg/errors"

var (
	ErrTemplateNotFound = errors.New("template is not found")
)

type Config struct {
	Dir string `yaml:"dir"`
}
