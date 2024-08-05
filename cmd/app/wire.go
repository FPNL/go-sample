//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/google/wire"

	"github.com/fpnl/go-sample/biz"
	"github.com/fpnl/go-sample/conf"
	"github.com/fpnl/go-sample/data"
	"github.com/fpnl/go-sample/pkg/tools"
	"github.com/fpnl/go-sample/server"
	"github.com/fpnl/go-sample/server/middleware"
	"github.com/fpnl/go-sample/service"
)

// initApp init kratos application.
func initApp(*conf.Project, *conf.Server, *conf.Data, *conf.Log, *slog.Logger) (*tools.App, func(), error) {
	panic(
		wire.Build(
			server.ProviderSet,
			data.ProviderSet,
			biz.ProviderSet,
			service.ProviderSet,
			middleware.ProviderSet,
			tools.NewApp,
		),
	)
}
