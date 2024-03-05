package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"sync"
	"time"
)

func NewFxPluginDescriptor(
	logger log.Logger,
	clientBuilderFactory smtp.ClientBuilderFactory,
	resolver smtp.MxResolver,
	emailParser smtp.EmailParser,
	dkimSignerFactory smtp.DkimSignerFactory,
) postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name:       "smtp",
			Kind:       postmanq.PluginKindSender,
			MinVersion: 1.0,
			Construct: func(ctx context.Context, pipeline postmanq.Pipeline, provider config.Provider) (postmanq.Plugin, error) {
				cfg := smtp.Config{
					Timeout: &smtp.TimeoutConfig{
						Connection: 5 * time.Minute,
						Hello:      5 * time.Minute,
						Mail:       5 * time.Minute,
						Rcpt:       5 * time.Minute,
						Data:       10 * time.Minute,
					},
				}
				err := provider.Populate(&cfg)
				if err != nil {
					return nil, err
				}

				builder, err := clientBuilderFactory.Create(ctx, cfg)
				if err != nil {
					return nil, err
				}

				ds, err := dkimSignerFactory.Create(ctx, cfg)
				if err != nil {
					return nil, err
				}

				p := &plugin{
					cfg:         cfg,
					logger:      logger.Named("smtp_plugin"),
					builder:     builder,
					resolver:    resolver,
					emailParser: emailParser,
					descriptors: collection.NewMap[string, *smtp.RecipientDescriptor](),
					noopTicker:  time.NewTicker(time.Minute),
					dkimSigner:  ds,
				}
				go p.startBackgroundProcess()
				return p, nil
			},
		},
	}
}

type plugin struct {
	cfg         smtp.Config
	logger      log.Logger
	builder     smtp.ClientBuilder
	resolver    smtp.MxResolver
	emailParser smtp.EmailParser
	descriptors collection.Map[string, *smtp.RecipientDescriptor]
	mtx         sync.Mutex
	noopTicker  *time.Ticker
	dkimSigner  smtp.DkimSigner
}

func (p *plugin) GetType() string {
	return "ActivityTypeSMTP"
}

func (p *plugin) OnEvent(ctx context.Context, event *postmanq.Event) (*postmanq.Event, error) {
	logger := p.logger.WithCtx(
		ctx,
		"uuid", event.Uuid,
		"from", event.From,
		"to", event.To,
	)
	logger.Debug("try to parse recipient email")
	email, err := p.emailParser.Parse(ctx, event.To)
	if err != nil {
		return nil, err
	}

	logger.Debug("recipient email parsed")
	p.mtx.Lock()
	descriptor, exists := p.descriptors.Get(email.Domain)
	if !exists {
		logger.Debug("try to resolve domain")
		mxRecords, err := p.resolver.Resolve(ctx, email.Domain)
		if err != nil {
			return nil, err
		}

		logger.Debug("domain resolved")
		descriptor = &smtp.RecipientDescriptor{
			Servers:    collection.NewSlice[*smtp.ServerDescriptor](),
			ModifiedAt: time.Now(),
		}
		for _, mxRecord := range mxRecords.Entries() {
			descriptor.Servers.Add(&smtp.ServerDescriptor{
				MxRecord:   mxRecord,
				Clients:    collection.NewSlice[smtp.Client](),
				ModifiedAt: time.Now(),
			})
		}

		p.descriptors.Set(email.Domain, descriptor)
	}
	p.mtx.Unlock()
	goto createClient

createClient:
	var cl smtp.Client
	var data []byte

	for _, server := range descriptor.Servers.Entries() {
		logger.Debugf("try to find exist client for %s", server.MxRecord.Host)
		for _, existsClient := range server.Clients.Entries() {
			if existsClient.HasStatus(smtp.ClientStatusBusy) {
				continue
			}

			cl = existsClient
			logger.Debug("exist client found")
			break
		}

		if cl == nil {
			logger.Debug("exist client not found")
			if server.HasMaxCountOfClients() {
				logger.Debug("max count of clients reached, waiting")
				continue
			}

			logger.Debug("try to create new client")
			cl, err = p.builder.Create(ctx, server.MxRecord.Host)
			if err != nil {
				logger.Error(err)
				server.SetMaxCountOfClientsOn()
				goto waitClient
			}

			server.Clients.Add(cl)
			logger.Debug("new client created")
		}
	}

	logger.Debugf("try to send HELO %s", p.cfg.Hostname)
	err = cl.Hello(ctx, p.cfg.Hostname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("HELO sent successfully")
	logger.Debugf("try to send MAIL %s", event.From)
	err = cl.Mail(ctx, event.From)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("MAIL sent successfully")
	logger.Debugf("try to send RCPT %s", event.From)
	err = cl.Rcpt(ctx, event.To)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("MAIL sent successfully")
	logger.Debug("try to sign data")
	data, err = p.dkimSigner.Sign(ctx, event.Data)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("data signed successfully")
	logger.Debug("try to send DATA")
	err = cl.Data(ctx, data)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("DATA sent successfully")
	return event, nil

waitClient:
	time.Sleep(time.Second * 10)
	goto createClient
}

func (p *plugin) startBackgroundProcess() {
	defer p.noopTicker.Stop()
	for {
		select {
		case <-p.noopTicker.C:
			for _, descriptor := range p.descriptors.Entries() {
				for _, server := range descriptor.Servers.Entries() {
					for _, cl := range server.Clients.Entries() {
						cl.Noop()
					}
				}
			}
		}
	}
}
