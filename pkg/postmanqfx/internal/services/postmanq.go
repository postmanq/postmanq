package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/reactivex/rxgo/v2"
)

type invoker struct {
	ctx               context.Context
	config            *postmanq.Config
	logger            log.Logger
	providerFactory   config.ProviderFactory
	pluginDescriptors []postmanq.PluginDescriptor
	Pipelines         map[string]*postmanq.Pipeline
}

func (i invoker) Configure() error {
	for _, pipelineCfg := range i.config.Pipelines {
		pipeline := &postmanq.Pipeline{
			Receivers:   make([]postmanq.Plugin, 0),
			Middlewares: make([]postmanq.Plugin, 0),
			Senders:     make([]postmanq.Plugin, 0),
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
					pipeline.Receivers = append(pipeline.Receivers, plugin)
				case descriptor.Kind&postmanq.PluginKindMiddleware == postmanq.PluginKindMiddleware:
					pipeline.Middlewares = append(pipeline.Middlewares, plugin)
				case descriptor.Kind&postmanq.PluginKindSender == postmanq.PluginKindSender:
					pipeline.Senders = append(pipeline.Senders, plugin)
				}
			}
		}

		i.Pipelines[pipelineCfg.Name] = pipeline
	}
	return nil
}

func (i invoker) Run(ctx context.Context) error {
	for _, pipeline := range i.Pipelines {
		childCtx, cancel := context.WithCancel(ctx)
		opts := []rxgo.Option{
			rxgo.WithContext(childCtx),
			rxgo.WithPool(i.config.PoolSize),
		}

		producers := make([]rxgo.Producer, len(pipeline.Receivers))
		for y, receiver := range pipeline.Receivers {
			producers[y] = func(ctx context.Context, next chan<- rxgo.Item) {
				err := receiver.OnReceive(ctx, next)
				if err != nil {
					next <- rxgo.Error(err)
				}
			}
		}

		nextObservable := rxgo.FromEventSource(rxgo.Defer(producers, opts...).Observe(), opts...)
		nextObservable.DoOnError(func(err error) {
			i.logger.Error(err)
			cancel()
		})
		if len(pipeline.Middlewares) > 0 {

		}

		for _, sender := range pipeline.Senders {
			nextObservable.DoOnNext(func(i interface{}) {
				err := sender.OnSend(childCtx, rxgo.Of(i))
				if err != nil {

				}
			})
		}
	}
	<-ctx.Done()
	return nil
}
