package main

import (
	"context"
	"github.com/postmanq/postmanq/pkg/commonfx"
	"github.com/postmanq/postmanq/pkg/commonfx/app"
	"github.com/postmanq/postmanq/pkg/postmanqfx"
	"github.com/postmanq/postmanq/pkg/postmanqfx/postmanq"
)

func main() {
	a := app.New(
		commonfx.Module,
		postmanqfx.Module,
	)
	a.Run(func(
		ctx context.Context,
		invoker postmanq.Invoker,
	) error {
		err := invoker.Configure(ctx)
		if err != nil {
			return err
		}

		return invoker.Run(ctx)
	})
}
