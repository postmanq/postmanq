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

func (a *App) Run() {
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

	var modules []interface{}
	for _, file := range files {
		plug, err := plugin.Open(fmt.Sprintf("%s/%s/module.so", args.ModuleDir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		symbol, err := plug.Lookup(module.Constructor)
		if err != nil {
			log.Fatal(err)
		}

		modules = append(modules, symbol)
	}

	app := fx.New(
		fx.Provide(func() Arguments {
			return args
		}),
		fx.Provide(modules...),
		fx.Invoke(func(params Params) {
			log.Println(params.Components)
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
