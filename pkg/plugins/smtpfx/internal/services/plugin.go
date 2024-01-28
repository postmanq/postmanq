package services

import (
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.temporal.io/sdk/workflow"
)

func NewFxPluginDescriptor() postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name:       "smtp",
			Kind:       postmanq.PluginKindSender,
			MinVersion: 1.0,
			Construct: func(provider config.Provider) (postmanq.Plugin, error) {
				return &plugin{}, nil
			},
		},
	}
}

type plugin struct {
}

func (p plugin) GetType() string {
	return "ActivityTypeSMTP"
}

func (p plugin) OnEvent(ctx workflow.Context, event *postmanqv1.Event) (*postmanqv1.Event, error) {
	//TODO implement me
	panic("implement me")
}
