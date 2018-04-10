package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
)

type AmlAlertRequest struct {
	BalanceID string `json:"balance_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

func (r *AmlAlertRequest) Populate(histRequest history.AmlAlertRequest) (error) {
	r.BalanceID = histRequest.BalanceID
	r.Amount = histRequest.Amount
	r.Reason = histRequest.Reason
	return nil
}
