package services

import (
	"fmt"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx/temporal"
	"go.temporal.io/sdk/client"
)

func NewFxClient(
	logger log.Logger,
	configProvider config.Provider,
) (temporal.Client, error) {
	cfg := new(temporal.Config)
	err := configProvider.PopulateByKey("temporal", cfg)
	if err != nil {
		return nil, err
	}

	opts := client.Options{
		HostPort: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Logger:   &logAdapter{logger.Named("tw_client")},
		ConnectionOptions: client.ConnectionOptions{
			MaxPayloadSize: 2097152 * 50,
		},
	}
	return client.Dial(opts)
}
