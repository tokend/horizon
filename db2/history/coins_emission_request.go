package history

import (
	"time"
)

// CoinsEmissionRequest  is a row of data from the `history_emission_requests` table
type CoinsEmissionRequest struct {
	RequestID    int64      `db:"request_id"`
	Receiver     string     `db:"receiver"`
	Reference    string     `db:"reference"`
	Issuer       string     `db:"issuer"`
	Amount       string     `db:"amount"`
	Asset        string     `db:"asset"`
	Approved     *bool      `db:"approved"`
	PreEmissions *string    `db:"preemissions"`
	Reason       *string    `db:"reason"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	ExchangeName string     `db:"exchange_name"` // comes from join on history_balances
}
