package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/tokend/regources"
)

// ASwapBidIndexAction renders a page of atomic swap bids
// filters by ownerID, base asset code
type ASwapBidShowAction struct {
	Action
	BidID  int64
	Record regources.AtomicSwapBid
}

// JSON is a method for actions.JSON
func (action *ASwapBidShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Record)
		},
	)
}

func (action *ASwapBidShowAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *ASwapBidShowAction) loadParams() {
	action.BidID = action.GetInt64("bid_id")
}

func (action *ASwapBidShowAction) loadRecords() {
	aswapBid, err := action.CoreQ().AtomicSwapBid().ByID(action.BidID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get aswap bid record")
		action.Err = &problem.ServerError
		return
	}

	if aswapBid == nil {
		action.Err = &problem.NotFound
		return
	}

	quoteAssets, err := action.CoreQ().AtomicSwapQuoteAsset().
		ByID([]int64{aswapBid.BidID}).Select()
	if err != nil {
		action.Log.WithError(err).Error(
			"Failed to get aswap bid quote assets")
		action.Err = &problem.ServerError
		return
	}
	var quoteAssetsSlice []regources.AssetPrice
	for _, quoteAsset := range quoteAssets {
		quoteAssetsSlice = append(quoteAssetsSlice, regources.AssetPrice{
			Asset: quoteAsset.QuoteAsset,
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	action.Record = resource.PopulateASwapBid(*aswapBid, quoteAssetsSlice)
}
