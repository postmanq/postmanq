package module

const (
	ConstructName = "Plugin"
)

type PluginConstruct func() Plugin

type Plugin struct {
	Constructs []interface{}
}
