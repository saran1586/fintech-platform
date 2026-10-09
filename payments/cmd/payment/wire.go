//go:build wireinject
// +build wireinject
//go:generate wire
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"
    "payments/internal/client"
	"payments/internal/biz"
	"payments/internal/conf"
	"payments/internal/data"
	"payments/internal/server"
	"payments/internal/service"
    "payments/internal/messaging"
	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

func provideAccountClientConfig(
	data *conf.Data,
) *conf.Data_AccountClient {
	return data.Account
}

func provideTransactionEndpoint(
	data *conf.Data,
) string {
	return data.TransactionClient.Endpoint
}


func wireApp(
	serverConfig *conf.Server,
	dataConfig *conf.Data,
	natsConfig *conf.NATS,
	logger *slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		client.ProviderSet,
		service.ProviderSet,
		provideAccountClientConfig,
		client.TransactionProviderSet,
		provideTransactionEndpoint,
		messaging.ProviderSet,
		newApp,
	))
}