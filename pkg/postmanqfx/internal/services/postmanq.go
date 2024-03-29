package services

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"go.uber.org/fx"
)

type InvokerParams struct {
	fx.In
	Logger             log.Logger
	ProviderFactory    config.ProviderFactory
	Provider           config.Provider
	PluginDescriptors  []postmanq.PluginDescriptor `group:"plugins"`
	WorkerFactory      temporal.WorkerFactory
	EventSenderFactory postmanq.EventSenderFactory
}

func NewFxInvoker(params InvokerParams) (postmanq.Invoker, error) {
	var configPipelines []postmanq.ConfigPipeline
	err := params.Provider.PopulateByKey("pipelines", &configPipelines)
	if err != nil {
		return nil, err
	}

	return &invoker{
		configPipelines:    configPipelines,
		logger:             params.Logger,
		providerFactory:    params.ProviderFactory,
		pluginDescriptors:  params.PluginDescriptors,
		workerFactory:      params.WorkerFactory,
		eventSenderFactory: params.EventSenderFactory,
		pipelines:          collection.NewMap[string, *postmanq.Pipeline](),
	}, nil
}

type invoker struct {
	configPipelines    []postmanq.ConfigPipeline
	logger             log.Logger
	providerFactory    config.ProviderFactory
	pluginDescriptors  []postmanq.PluginDescriptor
	workerFactory      temporal.WorkerFactory
	eventSenderFactory postmanq.EventSenderFactory
	pipelines          collection.Map[string, *postmanq.Pipeline]
}

func (i invoker) Configure(ctx context.Context) error {
	for _, configPipeline := range i.configPipelines {
		pipeline := &postmanq.Pipeline{
			Queue:       configPipeline.Queue,
			Receivers:   collection.NewSlice[postmanq.ReceiverPlugin](),
			Middlewares: collection.NewSlice[postmanq.WorkflowPlugin](),
			Senders:     collection.NewSlice[postmanq.WorkflowPlugin](),
		}

		for _, pluginCfg := range configPipeline.Plugins {
			for _, descriptor := range i.pluginDescriptors {
				if pluginCfg.Name != descriptor.Name {
					continue
				}

				provider, err := i.providerFactory.Create(config.Static(pluginCfg.Config))
				if err != nil {
					return err
				}

				plugin, err := descriptor.Construct(ctx, *pipeline, provider)
				if err != nil {
					return err
				}

				switch {
				case descriptor.Kind&postmanq.PluginKindReceiver == postmanq.PluginKindReceiver:
					pipeline.Receivers.Add(plugin.(postmanq.ReceiverPlugin))
				case descriptor.Kind&postmanq.PluginKindMiddleware == postmanq.PluginKindMiddleware:
					pipeline.Middlewares.Add(plugin.(postmanq.WorkflowPlugin))
				case descriptor.Kind&postmanq.PluginKindSender == postmanq.PluginKindSender:
					pipeline.Senders.Add(plugin.(postmanq.WorkflowPlugin))
				}
			}
		}

		i.pipelines.Set(configPipeline.Queue, pipeline)
	}
	return nil
}

func (i invoker) Run(ctx context.Context) error {
	for _, pipeline := range i.pipelines.Entries() {
		childCtx, cancel := context.WithCancel(ctx)
		for _, plugin := range pipeline.Receivers.Entries() {
			go func(ctx context.Context, cancel context.CancelFunc, plugin postmanq.ReceiverPlugin) {
				err := plugin.Receive(ctx)
				if err != nil {
					i.logger.Error(err)
					cancel()
				}
			}(childCtx, cancel, plugin)
		}

		workflowPlugins := collection.NewSlice[postmanq.WorkflowPlugin]()
		workflowPlugins.Add(pipeline.Middlewares.Entries()...)
		workflowPlugins.Add(pipeline.Senders.Entries()...)
		activityDescriptors := collection.NewSlice[temporal.ActivityDescriptor]()
		for _, plugin := range workflowPlugins.Entries() {
			activityDescriptors.Add(temporal.ActivityDescriptor{
				Type: plugin.GetType(),
				Func: plugin.OnEvent,
			})
		}

		sender := i.eventSenderFactory.Create(pipeline)
		worker, err := i.workerFactory.CreateByDescriptor(ctx, temporal.WorkerDescriptor{
			Workflow: temporal.WorkflowDescriptor{
				Type: fmt.Sprintf("WorkflowType%s", pipeline.Queue),
				Func: sender.SendEvent,
			},
			Activities: activityDescriptors.Entries(),
		})
		if err != nil {
			i.logger.Error(err)
			cancel()
		}

		err = worker.Run(temporal.InterruptCh())
		if err != nil {
			i.logger.Error(err)
			cancel()
		}
	}

	return nil
}
