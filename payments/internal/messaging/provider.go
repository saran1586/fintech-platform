package messaging

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewNATSConnection,
	NewPublisher,
	NewConsumer,
)