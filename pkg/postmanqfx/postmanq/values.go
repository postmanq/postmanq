package postmanq

import (
	"context"
	"github.com/postmanq/postmanq/pkg/collection"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/gen/postmanqv1"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
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
	Name        string
	Receivers   collection.Slice[ReceiverPlugin]
	Middlewares collection.Slice[WorkflowPlugin]
	Senders     collection.Slice[WorkflowPlugin]
}

type Config struct {
	PoolSize  int              `yaml:"pool_size"`
	Pipelines []ConfigPipeline `yaml:"pipelines"`
}

type ConfigPipeline struct {
	Name    string         `yaml:"name"`
	Plugins []ConfigPlugin `yaml:"plugins"`
}

type ConfigPlugin struct {
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config"`
}

type EventSenderFactoryParams struct {
	fx.In
	Ctx                     context.Context
	Logger                  log.Logger
	WorkflowExecutorFactory temporal.WorkflowExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
	ActivityExecutorFactory temporal.ActivityExecutorFactory[*postmanqv1.Event, *postmanqv1.Event]
}
