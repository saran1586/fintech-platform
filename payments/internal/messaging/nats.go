package messaging

import (
	"github.com/nats-io/nats.go"
	"payments/internal/conf"
)

func NewNATSConnection(config *conf.NATS) (*nats.Conn, error) {
	return nats.Connect(config.Addr)
}