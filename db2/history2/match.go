package history2

import (
	"gitlab.com/tokend/go/xdr"
	"time"
)

// Match is a row of data from the `matches` table
type Match struct {
	ID string `db:"id"`

	OrderBookID   uint64     `db:"order_book_id"`
	OperationID   int64      `db:"operation_id"`
	ParticipantID string     `db:"participant_id"`
	BaseAmount    int64      `db:"base_amount"`
	QuoteAmount   int64      `db:"quote_amount"`
	BaseAsset     string     `db:"base_asset"`
	QuoteAsset    string     `db:"quote_asset"`
	Price         int64      `db:"price"`
	CreatedAt     *time.Time `db:"created_at"`
}

// NewMatch returns new instance of Match
func NewMatch(
	base, quote xdr.AssetCode,
	orderBookID xdr.Uint64,
	operationID int64,
	atom xdr.ClaimOfferAtom,
) Match {
	return Match{
		ParticipantID: atom.BAccountId.Address(),
		OrderBookID:   uint64(orderBookID),
		OperationID:   operationID,
		BaseAmount:    int64(atom.BaseAmount),
		QuoteAmount:   int64(atom.QuoteAmount),
		BaseAsset:     string(base),
		QuoteAsset:    string(quote),
		Price:         int64(atom.CurrentPrice),
	}
}
