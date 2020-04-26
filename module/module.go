package module

const (
	ConstructName = "PqModule"
)

type DescriptorConstruct func() Descriptor

type Descriptor struct {
	Constructs []interface{}
}
