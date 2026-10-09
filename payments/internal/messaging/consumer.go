package messaging

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"payments/internal/event"
)

type Consumer struct {
	conn   *nats.Conn
	logger *slog.Logger
}

func NewConsumer(conn *nats.Conn,logger *slog.Logger,) *Consumer {
	return &Consumer{
		conn:   conn,
		logger: logger,
	}
}

func (c *Consumer) Run() error {
	return c.SubscribePaymentInitiated()
}

func (c *Consumer) SubscribePaymentInitiated() error {
	_, err := c.conn.Subscribe(
		"payments.initiated.*",
		func(msg *nats.Msg) {
			var paymentEvent event.PaymentInitiatedEvent

			if err := json.Unmarshal(msg.Data, &paymentEvent); err != nil {
				c.logger.Error(
					"failed to unmarshal payment initiated event",
					"error",
					err,
				)
				return
			}

			c.logger.Info(
				"payment initiated event received",
				"payment_id",
				paymentEvent.PaymentID,
				"from_account_id",
				paymentEvent.FromAccountID,
				"to_account_id",
				paymentEvent.ToAccountID,
				"amount",
				paymentEvent.Amount,
				"currency",
				paymentEvent.Currency,
				"reference_number",
				paymentEvent.ReferenceNumber,
			)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) Stop() error {
	return c.conn.Drain()
}