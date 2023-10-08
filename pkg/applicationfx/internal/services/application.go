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

func New() application.Application {
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

	var modules []interface{}
	for _, file := range files {
		moduleName := fmt.Sprintf("%s/%s.so", args.ModuleDir, file.Name())
		plug, err := plugin.Open(moduleName)
		if err != nil {
			log.Fatal(err)
		}

		symbol, err := plug.Lookup(application.ModuleSymName)
		if err != nil {
			log.Fatal(err)
		}

		module, ok := symbol.(*fx.Option)
		if !ok {
			log.Fatal(fmt.Errorf("can`t cast symbol=%T to module.PluginConstruct in mudule %s", symbol, moduleName))
		}

		modules = append(modules, *module)
	}

	ctx, _ := context.WithCancel(context.Background())
	fxApp := fx.New(
		fx.Provide(func() application.Arguments {
			return args
		}),
		fx.Provide(func() context.Context {
			return ctx
		}),
		fx.Provide(modules...),
		fx.Invoke(invoker),
	)

	fxApp.Run()
}
