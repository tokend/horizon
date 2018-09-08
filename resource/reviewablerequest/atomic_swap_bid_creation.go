package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateASwapBidCreationRequest(histRequest history.AtomicSwapBidCreation) (
	*regources.AtomicSwapBidCreation, error,
) {
	return &regources.AtomicSwapBidCreation{
		BaseBalance: histRequest.BaseBalance,
		BaseAmount:  regources.Amount(histRequest.BaseAmount),
		QuoteAssets: histRequest.QuoteAssets,
		Details:     histRequest.Details,
	}, nil
}
