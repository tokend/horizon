package reviewablerequest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
	"time"
)

type SaleCreationRequest struct {
	BaseAsset         string                   `json:"base_asset"`
	DefaultQuoteAsset string                   `json:"default_quote_asset"`
	StartTime         time.Time                `json:"start_time"`
	EndTime           time.Time                `json:"end_time"`
	SoftCap           string                   `json:"soft_cap"`
	HardCap           string                   `json:"hard_cap"`
	Details           map[string]interface{}   `json:"details"`
	QuoteAssets       []history.SaleQuoteAsset `json:"quote_assets"`
}

func (r *SaleCreationRequest) Populate(histRequest history.SaleRequest) {
	r.BaseAsset = histRequest.BaseAsset
	r.DefaultQuoteAsset = histRequest.DefaultQuoteAsset
	r.StartTime = histRequest.StartTime
	r.EndTime = histRequest.EndTime
	r.SoftCap = histRequest.SoftCap
	r.HardCap = histRequest.HardCap
	r.Details = histRequest.Details
	r.QuoteAssets = histRequest.QuoteAssets
}

func (r *SaleCreationRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.SaleRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.SaleRequest")
	}

	r.Populate(histRequest)
	return nil
}
