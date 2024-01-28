package smtp

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
)

type MxResolver interface {
	Resolve(ctx context.Context, domain string) (collection.Slice[MxRecord], error)
}
