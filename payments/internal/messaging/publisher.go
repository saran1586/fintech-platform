package messaging

import (
	"encoding/json"
	"strconv"

	"github.com/nats-io/nats.go"
	"payments/internal/event"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{
		conn: conn,
	}
}

func (p *Publisher) PublishPaymentInitiated(e event.PaymentInitiatedEvent) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	subject := "payments.initiated." + strconv.FormatInt(e.PaymentID, 10)

	return p.conn.Publish(subject, data)
}