package cli

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/postmanq/postmanq/module"
	"go.uber.org/fx"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

type App struct {
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
		args.ModuleDir = fmt.Sprintf("%s/module", filepath.Dir(ex))
	}

	files, err := ioutil.ReadDir(args.ModuleDir)
	if err != nil {
		log.Fatal(err)
	}

	var constructs []interface{}
	for _, file := range files {
		moduleName := fmt.Sprintf("%s/%s/module.so", args.ModuleDir, file.Name())
		plug, err := plugin.Open(moduleName)
		if err != nil {
			log.Fatal(err)
		}

		symbol, err := plug.Lookup(module.ConstructName)
		if err != nil {
			log.Fatal(err)
		}

		descriptorConstruct, ok := symbol.(*module.PluginConstruct)
		if !ok {
			log.Fatal(fmt.Errorf("can`t cast symbol=%T to module.PluginConstruct in mudule %s", symbol, moduleName))
		}

		descriptor := (*descriptorConstruct)()
		constructs = append(constructs, descriptor.Constructs...)
	}

	app := fx.New(
		fx.Provide(func() Arguments {
			return args
		}),
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		fx.Provide(constructs...),
		fx.Invoke(invoker),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
