package services

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/postmanq/postmanq/pkg/applicationfx/application"
	"go.uber.org/fx"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

func NewFxApplication() application.Application {
	return &app{}
}

type app struct {
}

func (a *app) Run(invoker interface{}) {
	var args application.Arguments
	_, err := flags.Parse(&args)
	if err != nil {
		return
	}

	if len(args.ModuleDir) == 0 {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		args.ModuleDir = fmt.Sprintf("%s/plugins", filepath.Dir(ex))
	}

	files, err := ioutil.ReadDir(args.ModuleDir)
	if err != nil {
		log.Fatal(err)
	}

	var constructs []interface{}
	for _, file := range files {
		moduleName := fmt.Sprintf("%s/%s.so", args.ModuleDir, file.Name())
		plug, err := plugin.Open(moduleName)
		if err != nil {
			log.Fatal(err)
		}

		symbol, err := plug.Lookup(application.ConstructName)
		if err != nil {
			log.Fatal(err)
		}

		_, ok := symbol.(*application.PluginConstruct)
		if !ok {
			log.Fatal(fmt.Errorf("can`t cast symbol=%T to module.PluginConstruct in mudule %s", symbol, moduleName))
		}

		//descriptor := (*descriptorConstruct)()
		//constructs = append(constructs, descriptor.Constructs...)
	}

	ctx := context.Background()
	fxApp := fx.New(
		fx.Provide(func() application.Arguments {
			return args
		}),
		fx.Provide(func() context.Context {
			return ctx
		}),
		fx.Provide(constructs...),
		fx.Invoke(invoker),
	)

	if err := fxApp.Start(ctx); err != nil {
		log.Fatal(err)
	}

	if err := fxApp.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}
