package plugin

const (
	ConstructName = "Plugin"
)

type Construct func() Plugin

type Plugin struct {
	Constructs []interface{}
}
