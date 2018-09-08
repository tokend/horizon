package resource

import (
	"strconv"

	"encoding/json"
	"time"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/regources"
)

func PopulateASwapBid(
	bid core.AtomicSwapBidEntry,
	quoteAssets []regources.AssetPrice,
) regources.AtomicSwapBid {
	var details map[string]interface{}
	_ = json.Unmarshal([]byte(bid.Details), &details)

	return regources.AtomicSwapBid{
		ID:              strconv.FormatInt(bid.BidID, 10),
		PT:              strconv.FormatInt(bid.BidID, 10),
		OwnerID:         bid.OwnerID,
		BaseAsset:       bid.BaseAsset,
		BaseBalanceID:   bid.BaseBalanceID,
		AvailableAmount: regources.Amount(bid.AvailableAmount),
		LockedAmount:    regources.Amount(bid.LockedAmount),
		CreatedAt:       time.Unix(bid.CreatedAt, 0).UTC(),
		IsCanceled:      bid.IsCanceled,
		Details:         details,
		QuoteAssets:     quoteAssets,
	}
}
