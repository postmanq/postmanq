package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out

	Descriptor PluginDescriptor `group:"plugins"`
}

type Params struct {
	fx.In
	Ctx               context.Context
	Config            *Config
	Logger            log.Logger
	ProviderFactory   config.ProviderFactory
	PluginDescriptors []PluginDescriptor `group:"plugins"`
}

type PluginDescriptor struct {
	Name       string
	Kind       PluginKind
	Construct  PluginConstruct
	MinVersion float32
}

type PluginKind int

const (
	PluginKindUnknown    PluginKind = 0
	PluginKindReceiver   PluginKind = 1
	PluginKindSender     PluginKind = 2
	PluginKindMiddleware PluginKind = 4
)

type Pipeline struct {
	Queue       string
	Receivers   collection.Slice[ReceiverPlugin]
	Middlewares collection.Slice[WorkflowPlugin]
	Senders     collection.Slice[WorkflowPlugin]
}

type Config struct {
	Pipelines []ConfigPipeline `yaml:"pipelines"`
}

type ConfigPipeline struct {
	Queue   string         `yaml:"queue"`
	Plugins []ConfigPlugin `yaml:"plugins"`
}

type ConfigPlugin struct {
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config"`
}
