package reviewablerequest

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateAtomicSwapRequest(histRequest history.AtomicSwap) (
	*regources.AtomicSwap, error,
) {
	return &regources.AtomicSwap{
		BidID:      strconv.FormatUint(histRequest.BidID, 10),
		BaseAmount: regources.Amount(histRequest.BaseAmount),
		QuoteAsset: histRequest.QuoteAsset,
	}, nil
}
