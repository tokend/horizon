package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type IssuanceRequest struct {
	Asset    string `json:"asset"`
	Amount   string `json:"amount"`
	Receiver string `json:"receiver"`
}

func (r *IssuanceRequest) Populate(histRequest history.IssuanceRequest) {
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Receiver = histRequest.Receiver
}

func (r *IssuanceRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.IssuanceRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.IssuanceRequest")
	}

	r.Populate(histRequest)
	return nil
}
