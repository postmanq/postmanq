package app

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"go.uber.org/fx"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

type Arguments struct {
	ConfigFilename string `short:"c" long:"config" description:"A path to config file" required:"true"`
	ModuleDir      string `short:"d" long:"dir" description:"A directory contains postmanq modules"`
}

const (
	ModuleSymName = "Module"
)

func New(modules ...fx.Option) *App {
	return &App{
		modules: modules,
	}
}

type App struct {
	modules []fx.Option
}

func (a *App) Run(invoker interface{}) {
	var args Arguments
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

	for _, file := range files {
		moduleName := fmt.Sprintf("%s/%s", args.ModuleDir, file.Name())
		plug, err := plugin.Open(moduleName)
		if err != nil {
			log.Fatal(err)
		}

		symbol, err := plug.Lookup(ModuleSymName)
		if err != nil {
			log.Fatal(err)
		}

		module, ok := symbol.(*fx.Option)
		if !ok {
			log.Fatal(fmt.Errorf("can`t cast symbol=%T to module.PluginConstruct in mudule %s", symbol, moduleName))
		}

		a.modules = append(a.modules, *module)
	}

	ctx, _ := context.WithCancel(context.Background())
	opts := append(
		[]fx.Option{},
		fx.Provide(func() Arguments {
			return args
		}),
		fx.Provide(func() context.Context {
			return ctx
		}),
		fx.Invoke(invoker),
	)
	opts = append(opts, a.modules...)

	fxApp := fx.New(opts...)
	fxApp.Run()
}
