package component

import (
	"github.com/byorty/postmanq_sender/plugin"
)

type CompleteSender struct {
	*SimpleSender
}

func NewCompleteSender(sender plugin.Sender) *CompleteSender {
	return &CompleteSender{
		SimpleSender: NewSimpleSender(sender, nil),
	}
}

func (c *CompleteSender) Run(numCPU int) {
	for i := 0; i < numCPU; i++ {
		go func() {
			for s := range c.sendings {
				err := c.sender.OnSend(s)

				if err == nil {
					s.Complete()
				} else {
					s.Abort(err)
				}
			}
		}()
	}
}
