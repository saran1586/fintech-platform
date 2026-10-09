package client

import (
	"account/internal/conf"

	"github.com/google/wire"
)

func NewCustomerClientConfig(data *conf.Data) *conf.CustomerClient {
	return data.CustomerClient
}

var ProviderSet = wire.NewSet(
	NewCustomerClientConfig,
	NewCustomerClient,
)