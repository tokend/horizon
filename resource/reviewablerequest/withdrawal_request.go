package reviewablerequest

import (
	"encoding/json"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type WithdrawalRequest struct {
	BalanceID       string `json:"balance_id"`
	Amount          string `json:"amount"`
	FixedFee        string `json:"fixed_fee"`
	PercentFee      string `json:"percent_fee"`
	ExternalDetails string `json:"external_details"`
	DestAssetCode   string `json:"dest_asset_code"`
	DestAssetAmount string `json:"dest_asset_amount"`
	ReviewerDetails map[string]interface{} `json:"reviewer_details"`
}

func (r *WithdrawalRequest) Populate(histRequest history.WithdrawalRequest) {
	r.BalanceID = histRequest.BalanceID
	r.Amount = histRequest.Amount
	r.FixedFee = histRequest.FixedFee
	r.PercentFee = histRequest.PercentFee
	r.ExternalDetails = histRequest.ExternalDetails
	r.DestAssetCode = histRequest.DestAssetCode
	r.DestAssetAmount = histRequest.DestAssetAmount
	r.ReviewerDetails = histRequest.ReviewerDetails
}

func (r *WithdrawalRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.WithdrawalRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.WithdrawalRequest")
	}

	r.Populate(histRequest)
	return nil
}
