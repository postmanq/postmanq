package temporal

import (
	"github.com/postmanq/postmanq/pkg/configfx/config"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"github.com/postmanq/postmanq/pkg/temporalfx/temporal"
	"go.temporal.io/sdk/client"
)

func NewFxClient(
	logger log.Logger,
	configProvider config.Provider,
) (temporal.Client, error) {
	cfg := new(temporal.Config)
	err := configProvider.PopulateByKey("temporal.client", cfg)
	if err != nil {
		return nil, err
	}

	opts := client.Options{
		HostPort: cfg.Address,
		Logger:   &logAdapter{logger.Named("tw_client")},
		ConnectionOptions: client.ConnectionOptions{
			MaxPayloadSize: 2097152 * 50,
		},
	}
	return client.Dial(opts)
}
