package services_test

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config_mock"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx/log_mock"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/internal/services"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp_mocks"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPluginTestSuite(t *testing.T) {
	suite.Run(t, new(PluginTestSuite))
}

type PluginTestSuite struct {
	testutils.Suite
	plugin        postmanq.WorkflowPlugin
	clientBuilder *smtp_mocks.MockClientBuilder
	parser        *smtp_mocks.MockEmailParser
	resolver      *smtp_mocks.MockMxResolver
	domains       collection.Slice[string]
	rand          *rand.Rand
	mtx           sync.Mutex
}

func (s *PluginTestSuite) SetupSuite() {
	s.Suite.SetupSuite()

	provider := config_mock.NewMockProvider(s.Ctrl)
	provider.EXPECT().Populate(gomock.Any()).Return(nil)

	s.clientBuilder = smtp_mocks.NewMockClientBuilder(s.Ctrl)
	clientBuilderFactory := smtp_mocks.NewMockClientBuilderFactory(s.Ctrl)
	clientBuilderFactory.EXPECT().Create(s.Ctx, gomock.Any()).Return(s.clientBuilder, nil)

	dkimSigner := smtp_mocks.NewMockDkimSigner(s.Ctrl)
	dkimSigner.EXPECT().Sign(s.Ctx, gomock.Any()).Return([]byte{}, nil).AnyTimes()
	dkimSignerFactory := smtp_mocks.NewMockDkimSignerFactory(s.Ctrl)
	dkimSignerFactory.EXPECT().Create(s.Ctx, gomock.Any()).Return(dkimSigner, nil)

	logger := log_mock.NewMockLogger(s.Ctrl)
	logger.EXPECT().Named(gomock.Any()).Return(logger)
	logger.EXPECT().WithCtx(
		gomock.Any(),
		"uuid", gomock.Any(),
		"from", gomock.Any(),
		"to", gomock.Any()).Return(logger).AnyTimes()
	logger.EXPECT().Debug(gomock.Any()).Return().AnyTimes()
	logger.EXPECT().Debugf(gomock.Any(), gomock.Any()).Return().AnyTimes()

	s.resolver = smtp_mocks.NewMockMxResolver(s.Ctrl)
	s.parser = smtp_mocks.NewMockEmailParser(s.Ctrl)
	res := services.NewFxPluginDescriptor(
		logger,
		clientBuilderFactory,
		s.resolver,
		s.parser,
		dkimSignerFactory,
	)
	plugin, err := res.Descriptor.Construct(s.Ctx, postmanq.Pipeline{}, provider)
	s.Nil(err)

	s.plugin = plugin.(postmanq.WorkflowPlugin)
	s.rand = rand.New(rand.NewSource(time.Now().Unix()))
	s.domains = collection.ImportSlice[string](
		"gmail.com",
		"ya.ru",
		"mail.ru",
		"icloud.com",
		"outlook.com",
	)

	for _, domain := range s.domains.Entries() {
		domainMxRecords := collection.NewSlice[smtp.MxRecord]()
		for i := 0; i < 1+s.rand.Intn(3); i++ {
			host := fmt.Sprintf("%s.%s", randomdata.Alphanumeric(8), domain)
			domainMxRecords.Add(smtp.MxRecord{
				Host:      host,
				Priority:  uint16(i),
				CreatedAt: time.Now(),
			})

			clientMeth1 := func(_ context.Context, _ string) error {
				s.mtx.Lock()
				defer s.mtx.Unlock()
				time.Sleep(time.Second * time.Duration(s.rand.Intn(1)))
				return nil
			}
			client := smtp_mocks.NewMockClient(s.Ctrl)
			client.EXPECT().Hello(s.Ctx, gomock.Any()).DoAndReturn(clientMeth1).AnyTimes()
			client.EXPECT().Mail(s.Ctx, gomock.Any()).DoAndReturn(clientMeth1).AnyTimes()
			client.EXPECT().Rcpt(s.Ctx, gomock.Any()).DoAndReturn(clientMeth1).AnyTimes()
			client.EXPECT().Data(s.Ctx, gomock.Any()).DoAndReturn(func(_ context.Context, _ []byte) error {
				s.mtx.Lock()
				defer s.mtx.Unlock()
				time.Sleep(time.Second * time.Duration(s.rand.Intn(2)))
				return nil
			}).AnyTimes()
			client.EXPECT().HasStatus(gomock.Any()).DoAndReturn(func(status smtp.ClientStatus) bool {
				s.mtx.Lock()
				defer s.mtx.Unlock()
				if s.rand.Intn(5) == 3 {
					return false
				} else {
					return true
				}
			}).AnyTimes()
			s.clientBuilder.EXPECT().Create(s.Ctx, host).Return(client, nil).AnyTimes()
		}
		s.resolver.EXPECT().Resolve(s.Ctx, domain).Return(domainMxRecords, nil).AnyTimes()
	}
}

func (s *PluginTestSuite) TestOnEvent() {
	for i := 0; i < 100; i++ {
		s.Run(fmt.Sprintf("TestOnEvent#%d", i), func() {
			s.T().Parallel()
			s.mtx.Lock()
			localPart := randomdata.Alphanumeric(16)
			domain := s.domains.Get(s.rand.Intn(s.domains.Len()))
			email := fmt.Sprintf("%s@%s", localPart, domain)
			s.mtx.Unlock()

			s.parser.EXPECT().Parse(s.Ctx, email).Return(&smtp.EmailAddress{
				Address:   email,
				LocalPart: localPart,
				Domain:    domain,
			}, nil)

			_, err := s.plugin.OnEvent(s.Ctx, &postmanq.Event{
				Uuid: uuid.NewString(),
				To:   email,
			})
			s.Nil(err)
		})
	}
}
