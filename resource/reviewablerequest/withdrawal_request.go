package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
)

type WithdrawalRequest struct {
	BalanceID              string                 `json:"balance_id"`
	Amount                 string                 `json:"amount"`
	FixedFee               string                 `json:"fixed_fee"`
	PercentFee             string                 `json:"percent_fee"`
	PreConfirmationDetails map[string]interface{} `json:"pre_confirmation_details"`
	ExternalDetails        map[string]interface{} `json:"external_details"`
	DestAssetCode          string                 `json:"dest_asset_code"`
	DestAssetAmount        string                 `json:"dest_asset_amount"`
	ReviewerDetails        map[string]interface{} `json:"reviewer_details"`
}

func (r *WithdrawalRequest) Populate(histRequest history.WithdrawalRequest) error {
	r.BalanceID = histRequest.BalanceID
	r.Amount = histRequest.Amount
	r.FixedFee = histRequest.FixedFee
	r.PercentFee = histRequest.PercentFee
	r.ExternalDetails = histRequest.ExternalDetails
	r.DestAssetCode = histRequest.DestAssetCode
	r.DestAssetAmount = histRequest.DestAssetAmount
	r.ReviewerDetails = histRequest.ReviewerDetails
	r.PreConfirmationDetails = histRequest.PreConfirmationDetails
	return nil
}
