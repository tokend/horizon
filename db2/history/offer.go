package history

import (
	"time"

	"gitlab.com/tokend/horizon/db2"
)

type Offer struct {
	db2.TotalOrderID
	OfferID           int64     `db:"offer_id"`
	BaseAsset         string    `db:"base_asset"`
	QuoteAsset        string    `db:"quote_asset"`
	InitialBaseAmount int64     `db:"initial_base_amount"`
	CurrentBaseAmount int64     `db:"current_base_amount"`
	Price             int64     `db:"price"`
	OwnerID           string    `db:"owner_id"`
	IsCanceled        bool      `db:"is_canceled"`
	CreatedAt         time.Time `db:"created_at"`
}
