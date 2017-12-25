package reviewablerequest

import (
	"time"
	"encoding/json"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SaleCreationRequest struct {
	BaseAsset  string                 `json:"base_asset"`
	QuoteAsset string                 `json:"quote_asset"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Price      string                 `json:"price"`
	SoftCap    string                 `json:"soft_cap"`
	HardCap    string                 `json:"hard_cap"`
	Details    map[string]interface{} `json:"details"`
}

func (r *SaleCreationRequest) Populate(histRequest history.SaleRequest) {
	r.BaseAsset = histRequest.BaseAsset
	r.QuoteAsset = histRequest.QuoteAsset
	r.StartTime = histRequest.StartTime
	r.EndTime = histRequest.EndTime
	r.Price = histRequest.Price
	r.SoftCap = histRequest.SoftCap
	r.HardCap = histRequest.HardCap
	r.Details = histRequest.Details
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

