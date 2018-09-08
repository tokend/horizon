package reviewablerequest

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateASwapRequest(histRequest history.AtomicSwap) (
	*regources.AtomicSwap, error,
) {
	return &regources.AtomicSwap{
		BidID:      strconv.FormatUint(histRequest.BidID, 10),
		BaseAmount: regources.Amount(histRequest.BaseAmount),
		QuoteAsset: histRequest.QuoteAsset,
	}, nil
}
