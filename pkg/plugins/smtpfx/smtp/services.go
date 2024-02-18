package smtp

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
)

type MxResolver interface {
	Resolve(ctx context.Context, domain string) (collection.Slice[MxRecord], error)
}

type ClientBuilderFactory interface {
	Create(ctx context.Context, cfg Config) (ClientBuilder, error)
}

type ClientBuilder interface {
	Create(ctx context.Context, hostname string) (Client, error)
}

type Client interface {
	Hello(ctx context.Context, localName string) error
	Mail(ctx context.Context, from string) error
	Rcpt(ctx context.Context, to string) error
	Data(ctx context.Context, data []byte) error
	HasStatus(ClientStatus) bool
	Noop() error
}

type EmailParser interface {
	Parse(email string) (*EmailAddress, error)
}

type DkimSignerFactory interface {
	Create(cfg Config) (DkimSigner, error)
}

type DkimSigner interface {
	Sign(data []byte) ([]byte, error)
}
