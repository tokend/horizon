package history2

import (
	"gitlab.com/tokend/go/xdr"
)

// Match is a row of data from the `matches` table
type Match struct {
	OrderBookID   uint64 `db:"order_book_id"`
	ParticipantID string `db:"participant_id"`
	BaseAmount    int64  `db:"base_amount"`
	QuoteAmount   int64  `db:"quote_amount"`
	BaseAsset     string `db:"base_asset"`
	QuoteAsset    string `db:"quote_asset"`
	Price         int64  `db:"price"`
}

func NewMatch(base, quote xdr.AssetCode, orderBookID xdr.Uint64, atom xdr.ClaimOfferAtom) Match {
	return Match{
		OrderBookID:   uint64(orderBookID),
		BaseAmount:    int64(atom.BaseAmount),
		QuoteAmount:   int64(atom.QuoteAmount),
		BaseAsset:     string(base),
		QuoteAsset:    string(quote),
		ParticipantID: atom.BAccountId.Address(),
		Price:         int64(atom.CurrentPrice),
	}
}
