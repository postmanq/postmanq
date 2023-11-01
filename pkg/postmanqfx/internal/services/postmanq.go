package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/collection"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
)

type invoker struct {
	ctx               context.Context
	config            *postmanq.Config
	logger            log.Logger
	providerFactory   config.ProviderFactory
	pluginDescriptors []postmanq.PluginDescriptor
	workerFactory     temporal.WorkerFactory
	pipelines         map[string]*postmanq.Pipeline
}

func (i invoker) Configure() error {
	for _, pipelineCfg := range i.config.Pipelines {
		pipeline := &postmanq.Pipeline{
			Name:        pipelineCfg.Name,
			Receivers:   collection.NewSlice[postmanq.ReceiverPlugin](),
			Middlewares: collection.NewSlice[postmanq.WorkflowPlugin](),
			Senders:     collection.NewSlice[postmanq.WorkflowPlugin](),
		}

		for _, pluginCfg := range pipelineCfg.Plugins {
			for _, descriptor := range i.pluginDescriptors {
				if pluginCfg.Name != descriptor.Name {
					continue
				}

				provider, err := i.providerFactory.Create(config.Static(pluginCfg.Config))
				if err != nil {
					return err
				}

				plugin, err := descriptor.Construct(provider)
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

		i.pipelines[pipelineCfg.Name] = pipeline
	}
	return nil
}

func (i invoker) Run(ctx context.Context) error {
	for _, pipeline := range i.pipelines {
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

		worker, err := i.workerFactory.CreateByDescriptor(ctx, temporal.WorkerDescriptor{
			Workflow:   nil,
			Activities: nil,
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
	<-ctx.Done()
	return nil
}
