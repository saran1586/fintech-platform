 package event

type PaymentInitiatedEvent struct {
	PaymentID       int64  `json:"payment_id"`
	FromAccountID   int64  `json:"from_account_id"`
	ToAccountID     int64  `json:"to_account_id"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	ReferenceNumber string `json:"reference_number"`
}