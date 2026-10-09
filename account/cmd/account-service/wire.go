//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"account/internal/biz"
	"account/internal/client"
	"account/internal/conf"
	"account/internal/data"
	"account/internal/server"
	"account/internal/service"
	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet , client.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
