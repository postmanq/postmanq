package application

type Arguments struct {
	ConfigFilename string `short:"c" long:"config" description:"A path to config file" required:"true"`
	ModuleDir      string `short:"d" long:"dir" description:"A directory contains postmanq modules"`
}

type PluginKind int

const (
	ConstructName = ""

	UnknownKind    PluginKind = 0
	ReceiverKind   PluginKind = 1
	SenderKind     PluginKind = 2
	MiddlewareKind PluginKind = 4
)
