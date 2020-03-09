package component

type Pipe struct {
	components []interface{}
}

func NewPipe(components []interface{}) *Pipe {
	return &Pipe{
		components: components,
	}
}

func (c *Pipe) Bootstrap() {

}
