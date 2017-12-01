package reviewablerequest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
)

type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

func (r *PreIssuanceRequest) Populate(histRequest history.PreIssuanceRequest) {
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Signature = histRequest.Signature
	r.Reference = histRequest.Reference
}

func (r *PreIssuanceRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.PreIssuanceRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.PreIssuanceRequest")
	}

	r.Populate(histRequest)
	return nil
}
