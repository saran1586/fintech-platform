package client

import (
	"payments/internal/biz"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewAccountClient,
	wire.Bind(new(biz.AccountClient), new(*AccountClient)),
)
var TransactionProviderSet = wire.NewSet(
	NewTransactionClient,
	wire.Bind(new(biz.TransactionClient), new(*TransactionClient)),
)