package application

type Arguments struct {
	ConfigFilename string `short:"c" long:"config" description:"A path to config file" required:"true"`
	ModuleDir      string `short:"d" long:"dir" description:"A directory contains postmanq modules"`
}

const (
	ModuleSymName = "Module"
)
