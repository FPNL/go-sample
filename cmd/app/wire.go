//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/google/wire"

	"oltp/biz"
	"oltp/conf"
	"oltp/data"
	"oltp/server"
	"oltp/server/middleware"
	"oltp/service"
)

// initApp init kratos application.
func initApp(*conf.Project, *conf.Server, *conf.Data, *slog.Logger) (*app, func(), error) {
	panic(
		wire.Build(
			server.ProviderSet,
			data.ProviderSet,
			biz.ProviderSet,
			service.ProviderSet,
			middleware.ProviderSet,
			newApp,
		),
	)
}
