package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/flosch/pongo2/v6"
	"github.com/postmanq/postmanq/pkg/commonfx/collection"
	"github.com/postmanq/postmanq/pkg/commonfx/configfx/config"
	"github.com/postmanq/postmanq/pkg/plugins/templatefx/template"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
	"os"
	"path/filepath"
)

func NewFxPluginDescriptor() postmanq.Result {
	return postmanq.Result{
		Descriptor: postmanq.PluginDescriptor{
			Name:       "template",
			Kind:       postmanq.PluginKindMiddleware,
			MinVersion: 1.0,
			Construct: func(ctx context.Context, pipeline postmanq.Pipeline, provider config.Provider) (postmanq.Plugin, error) {
				var cfg template.Config
				err := provider.Populate(&cfg)
				if err != nil {
					return nil, err
				}

				templates := collection.NewMap[string, *pongo2.Template]()
				err = filepath.Walk(cfg.Dir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}

					rel, err := filepath.Rel(cfg.Dir, path)
					if err != nil {
						return err
					}

					templates.Set(rel, pongo2.Must(pongo2.FromFile(path)))

					return nil
				})
				if err != nil {
					return nil, err
				}

				p := &plugin{
					templates: templates,
				}
				return p, nil
			},
		},
	}
}

type plugin struct {
	templates collection.Map[string, *pongo2.Template]
}

func (p *plugin) GetType() string {
	return "ActivityTypeTemplate"
}

func (p *plugin) OnEvent(ctx context.Context, event *postmanq.Event) (*postmanq.Event, error) {
	tpl, ok := p.templates.Get(event.Template)
	if !ok {
		return nil, template.ErrTemplateNotFound
	}

	decodedBuf, err := base64.URLEncoding.DecodeString(event.Vars)
	if err != nil {
		return nil, err
	}

	tplCtx := make(pongo2.Context)
	err = json.Unmarshal(decodedBuf, &tplCtx)
	if err != nil {
		return nil, err
	}

	out, err := tpl.Execute(tplCtx)
	if err != nil {
		return nil, err
	}

	event.Data = out
	return event, nil
}
