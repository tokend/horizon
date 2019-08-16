package resources

import (
	"fmt"
	"time"

	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

//NewAccount - creates new instance of account
func NewAtomicSwapAsk(core core2.AtomicSwapAsk) regources.AtomicSwapAsk {
	return regources.AtomicSwapAsk{
		Key: NewAtomicSwapAskKey(uint64(core.AskID)),
		Attributes: regources.AtomicSwapAskAttributes{
			AvailableAmount: regources.Amount(core.AvailableAmount),
			LockedAmount:    regources.Amount(core.LockedAmount),
			CreatedAt:       time.Unix(int64(core.CreatedAt), 0).UTC(),
			IsCanceled:      core.IsCanceled,
			Details:         core.Details,
		},
	}
}

//NewAccountKey - creates account key from address
func NewAtomicSwapAskKey(id uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprint(id),
		Type: regources.ATOMIC_SWAP_ASK,
	}
}

func NewAtomicSwapAskQuoteAsset(raw core2.AtomicSwapQuoteAsset) regources.QuoteAsset {
	return regources.QuoteAsset{
		Key: regources.Key{
			// TODO: Use artificial ID
			ID:   fmt.Sprintf("%s:%d", raw.QuoteAsset, raw.AskID),
			Type: regources.QUOTE_ASSETS,
		},
		Attributes: regources.QuoteAssetAttributes{
			Price: regources.Amount(raw.Price),
		},
	}
}
