package resources

import (
	"fmt"
	"time"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/generated"
)

//NewAccount - creates new instance of account
func NewAtomicSwapBid(core core2.AtomicSwapBid) regources.AtomicSwapBid {
	return regources.AtomicSwapBid{
		Key: regources.Key{
			ID:   fmt.Sprint(core.BidID),
			Type: regources.ATOMIC_SWAP_BID,
		},
		Attributes: regources.AtomicSwapBidAttributes{
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
		Type: regources.ATOMIC_SWAP_BID,
	}
}

func NewAtomicSwapBidQuoteAsset(raw core2.AtomicSwapQuoteAsset) regources.QuoteAsset {
	return regources.QuoteAsset{
		Key: regources.Key{
			ID:   raw.QuoteAsset,
			Type: regources.QUOTE_ASSETS,
		},
		Attributes: regources.QuoteAssetAttributes{
			Price: regources.Amount(raw.Price),
		},
	}
}
