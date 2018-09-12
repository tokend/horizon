package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/tokend/regources"
)

// ASwapBidIndexAction renders a page of atomic swap bids
// filters by ownerID, base asset code
type ASwapBidIndexAction struct {
	Action
	PagingParams db2.PageQuery
	BaseAsset    string
	OwnerID      string
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *ASwapBidIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *ASwapBidIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.BaseAsset = action.GetString("base_asset")
	action.OwnerID = action.GetString("owner_id")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"base_asset": action.BaseAsset,
		"owner_id":   action.OwnerID,
	}
}

func (action *ASwapBidIndexAction) loadRecords() {
	q := action.CoreQ().AtomicSwapBid()
	if action.BaseAsset != "" {
		q = q.ForBaseAsset([]string{action.BaseAsset})
	}
	if action.OwnerID != "" {
		q = q.ForOwner(action.OwnerID)
	}

	aswapBids, err := q.Page(action.PagingParams).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get aswap bid records")
		action.Err = &problem.ServerError
		return
	}

	if aswapBids == nil {
		action.Err = &problem.NotFound
		return
	}

	var bidIDs []int64
	for _, aswapBid := range aswapBids {
		bidIDs = append(bidIDs, aswapBid.BidID)
	}

	quoteAssets, err := action.CoreQ().AtomicSwapQuoteAsset().
		ByID(bidIDs).Select()

	quoteAssetsMap := map[int64][]regources.AssetPrice{}
	for _, quoteAsset := range quoteAssets {
		quoteAssetsMap[quoteAsset.BidID] =
			append(quoteAssetsMap[quoteAsset.BidID], regources.AssetPrice{
				Asset: quoteAsset.QuoteAsset,
				Price: regources.Amount(quoteAsset.Price),
			})
	}

	for _, aswapBid := range aswapBids {
		action.Page.Add(resource.PopulateASwapBid(aswapBid,
			quoteAssetsMap[aswapBid.BidID]))
	}
}

func (action *ASwapBidIndexAction) loadPage() {
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
