package main

import (
	"github.com/nkhang/pluto/internal/fx/projectapifx"
	"github.com/nkhang/pluto/internal/fx/toolapifx"
	"github.com/nkhang/pluto/pkg/fx/configfx"
	"github.com/nkhang/pluto/pkg/fx/dbfx"
	"github.com/nkhang/pluto/pkg/fx/ginfx"
	"github.com/nkhang/pluto/pkg/fx/loggerfx"
	"github.com/nkhang/pluto/pkg/fx/redisfx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		configfx.Initialize("pluto"),
		loggerfx.Invoke,
		dbfx.Module,
		redisfx.Module,
		toolapifx.Module,
		projectapifx.Module,
		ginfx.Module,
		fx.Invoke(initializer),
	).Run()
}
