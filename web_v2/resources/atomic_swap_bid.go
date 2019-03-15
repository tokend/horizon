package resources

import (
	"fmt"
	"time"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewAccount - creates new instance of account
func NewAtomicSwapBid(core core2.AtomicSwapBid) regources.AtomicSwapBid {
	return regources.AtomicSwapBid{
		Key: regources.Key{
			ID:   fmt.Sprint(core.BidID),
			Type: regources.TypeAtomicSwapBid,
		},
		Attributes: regources.AtomicSwapBidAttrs{
			AvailableAmount: regources.Amount(core.AvailableAmount),
			LockedAmount:    regources.Amount(core.LockedAmount),
			CreatedAt:       time.Unix(int64(core.CreatedAt), 0).UTC(),
			IsCanceled:      core.IsCanceled,
			Details:         core.Details,
		},
	}
}

//NewAccountKey - creates account key from address
func NewAtomicSwapBidKey(id uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprint(id),
		Type: regources.TypeAtomicSwapBid,
	}
}

func NewAtomicSwapBidQuoteAsset(raw core2.AtomicSwapQuoteAsset) regources.QuoteAsset {
	return regources.QuoteAsset{
		Key: regources.Key{
			ID:   raw.QuoteAsset,
			Type: regources.TypeQuoteAssets,
		},
		Attributes: regources.QuoteAssetAttrs{
			Price: regources.Amount(raw.Price),
		},
	}
}
