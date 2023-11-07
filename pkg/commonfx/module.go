package commonfx

import (
	"github.com/postmanq/postmanq/pkg/commonfx/configfx"
	"github.com/postmanq/postmanq/pkg/commonfx/logfx"
	"github.com/postmanq/postmanq/pkg/commonfx/temporalfx"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"commonfx",
		configfx.Module,
		logfx.Module,
		temporalfx.Module,
	)
)
