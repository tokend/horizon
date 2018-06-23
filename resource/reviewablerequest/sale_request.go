package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"time"
)

type SaleCreationRequest struct {
	BaseAsset           string                   `json:"base_asset"`
	DefaultQuoteAsset   string                   `json:"default_quote_asset"`
	StartTime           time.Time                `json:"start_time"`
	EndTime             time.Time                `json:"end_time"`
	SoftCap             string                   `json:"soft_cap"`
	HardCap             string                   `json:"hard_cap"`
	SaleType            base.Flag                `json:"sale_type"`
	BaseAssetForHardCap string                   `json:"base_asset_for_hard_cap"`
	Details             map[string]interface{}   `json:"details"`
	QuoteAssets         []history.SaleQuoteAsset `json:"quote_assets"`
	State               base.Flag                `json:"state"`
}

func (r *SaleCreationRequest) Populate(histRequest history.SaleRequest) error {
	r.BaseAsset = histRequest.BaseAsset
	r.DefaultQuoteAsset = histRequest.DefaultQuoteAsset
	r.StartTime = histRequest.StartTime
	r.EndTime = histRequest.EndTime
	r.SoftCap = histRequest.SoftCap
	r.HardCap = histRequest.HardCap
	r.Details = histRequest.Details
	r.QuoteAssets = histRequest.QuoteAssets
	r.SaleType.Value = int32(histRequest.SaleType)
	r.SaleType.Name = histRequest.SaleType.ShortString()
	r.BaseAssetForHardCap = histRequest.BaseAssetForHardCap
	r.State.Value = int32(histRequest.State)
	r.State.Name = histRequest.State.ShortString()
	return nil
}
