package services_test

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config_mock"
	"github.com/postmanq/postmanq/pkg/commonfx/testutils"
	"github.com/postmanq/postmanq/pkg/plugins/templatefx/internal/services"
	"github.com/postmanq/postmanq/pkg/plugins/templatefx/template"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"path"
	"testing"
)

var (
	expectedData = `Title

John Doe

product 1
100.000000

product 2
99.990000


successful message
`
)

func TestPluginTestSuite(t *testing.T) {
	suite.Run(t, new(PluginTestSuite))
}

type PluginTestSuite struct {
	testutils.Suite
	plugin postmanq.WorkflowPlugin
}

func (s *PluginTestSuite) SetupSuite() {
	s.Suite.SetupSuite()

	provider := config_mock.NewMockProvider(s.Ctrl)
	provider.EXPECT().Populate(gomock.Any()).Do(func(cfg *template.Config) {
		cfg.Dir = path.Join(s.Dir, "tpls")
		s.T().Log(cfg.Dir)
	})

	res := services.NewFxPluginDescriptor()
	plugin, err := res.Descriptor.Construct(s.Ctx, postmanq.Pipeline{}, provider)
	s.Nil(err)

	s.plugin = plugin.(postmanq.WorkflowPlugin)
}

func (s *PluginTestSuite) TestOnEvent() {
	event := &postmanq.Event{
		Template: randomdata.Alphanumeric(10),
	}
	_, err := s.plugin.OnEvent(s.Ctx, event)
	s.Equal(template.ErrTemplateNotFound, err)

	vars := map[string]any{
		"title": "Title",
		"user": map[string]any{
			"name": "John Doe",
			"purchases": []any{
				map[string]any{
					"name":  "product 1",
					"price": 100.00,
				},
				map[string]any{
					"name":  "product 2",
					"price": 99.99,
				},
			},
		},
		"is_true": true,
	}

	buf, err := json.Marshal(vars)
	s.Nil(err)

	encodedBuf := base64.URLEncoding.EncodeToString(buf)
	event = &postmanq.Event{
		Template: "index.html",
		Vars:     encodedBuf,
	}
	_, err = s.plugin.OnEvent(s.Ctx, event)
	s.Nil(err)
	s.Equal(expectedData, event.Data)
}
