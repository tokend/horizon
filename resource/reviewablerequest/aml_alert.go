package reviewablerequest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/history"
)

type AmlAlertRequest struct {
	BalanceID string `json:"balance_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

func (r *AmlAlertRequest) Populate(histRequest history.AmlAlertRequest) {
	r.BalanceID = histRequest.BalanceID
	r.Amount = histRequest.Amount
	r.Reason = histRequest.Reason
}

func (r *AmlAlertRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.AmlAlertRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.AmlAlertRequest")
	}

	r.Populate(histRequest)
	return nil
}
