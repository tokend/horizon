package history2

import (
	"gitlab.com/tokend/go/xdr"
)

// SaleParticipation is a row of data from the `account_specific_rules` table
type SaleParticipation struct {
	ID            uint64 `db:"id"`
	SaleID        uint64 `db:"sale_id"`
	ParticipantID string `db:"participant_id"`
	BaseAmount    int64  `db:"base_amount"`
	QuoteAmount   int64  `db:"quote_amount"`
	BaseAsset     string `db:"base_asset"`
	QuoteAsset    string `db:"quote_asset"`
	Price         int64  `db:"price"`
}

func NewSaleParticipation(base, quote string, saleID uint64, atom xdr.ClaimOfferAtom) SaleParticipation {
	return SaleParticipation{
		ID:            uint64(atom.OfferId),
		SaleID:        saleID,
		BaseAmount:    int64(atom.BaseAmount),
		QuoteAmount:   int64(atom.QuoteAmount),
		BaseAsset:     base,
		QuoteAsset:    quote,
		ParticipantID: atom.BAccountId.Address(),
		Price:         int64(atom.CurrentPrice),
	}
}
