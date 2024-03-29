package services

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"net"
	"strings"
	"time"
)

func NewFxMxResolver(
	logger log.Logger,
) smtp.MxResolver {
	return &mxResolver{
		logger: logger.Named("mx_resolver"),
	}
}

type mxResolver struct {
	logger log.Logger
}

func (m *mxResolver) Resolve(ctx context.Context, domain string) (collection.Slice[smtp.MxRecord], error) {
	logger := m.logger.WithCtx(ctx)
	logger.Debugf("trying to resolve mx records for domain %s", domain)
	mxes, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	logger.Debugf("got %d mx records", len(mxes))
	sl := collection.NewSlice[smtp.MxRecord](collection.WithSliceSize(len(mxes)))
	for i, mx := range mxes {
		mxRecord := smtp.MxRecord{
			Host:      strings.TrimRight(mx.Host, "."),
			Priority:  mx.Pref,
			CreatedAt: time.Now(),
		}
		sl.Set(i, mxRecord)
		logger.Debugf("created mx records for host %s", mxRecord.Host)
	}

	return sl, nil
}
