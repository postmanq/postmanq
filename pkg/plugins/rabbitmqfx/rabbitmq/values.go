package rabbitmq

type Config struct {
	Url   string `yaml:"url"`
	Queue string `yaml:"queue"`
}
