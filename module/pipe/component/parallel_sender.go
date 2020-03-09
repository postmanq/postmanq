package component

import (
	"github.com/byorty/postmanq_sender/plugin"
	"github.com/byorty/postmanq_sender/sending"
	"github.com/byorty/postmanq_sender/util/concurrent"
	"go.uber.org/multierr"
	"sync"
)

type ParallelSender struct {
	*SimpleSender
	senders []plugin.Sender
}

func NewParallelSender(senders []plugin.Sender, next SendingStep) *ParallelSender {
	return &ParallelSender{
		SimpleSender: NewSimpleSender(nil, next),
		senders:      senders,
	}
}

func (c *ParallelSender) Run(numCPU int) {
	numSenders := len(c.senders)

	for i := 0; i < numCPU; i++ {
		go func() {
			senderErrors := concurrent.NewSlice()
			wg := new(sync.WaitGroup)
			for s := range c.sendings {
				wg.Add(numSenders)

				for y := 0; y < numSenders; y++ {
					go func(sender plugin.Sender, send sending.Sending) {
						err := sender.OnSend(send)
						if err != nil {
							senderErrors.Append(err)
						}
						wg.Done()
					}(c.senders[y], s)
				}

				wg.Wait()

				if senderErrors.Len() == 0 {
					c.next.Send(s)
				} else {
					var err error
					for i := 0; i < senderErrors.Len(); i++ {
						err = multierr.Append(err, senderErrors.Get(i).(error))
					}

					s.Abort(err)
				}

				senderErrors.Clear()
			}
		}()
	}
}
