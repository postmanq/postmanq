package main

import (
	"github.com/postmanq/postmanq/cli"
	pc "github.com/postmanq/postmanq/module/pipe/component"
)

func main() {
	app := new(cli.App)
	app.Run(func(runner *pc.Runner) error {
		return runner.Run()
	})
}
