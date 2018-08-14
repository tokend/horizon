package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateWithdrawalRequest(histRequest history.WithdrawalRequest) (
	r *reviewablerequest2.WithdrawalRequest, err error,
) {
	r = &reviewablerequest2.WithdrawalRequest{}
	r.BalanceID = histRequest.BalanceID
	r.Amount = histRequest.Amount
	r.FixedFee = histRequest.FixedFee
	r.PercentFee = histRequest.PercentFee
	r.ExternalDetails = histRequest.ExternalDetails
	r.DestAssetCode = histRequest.DestAssetCode
	r.DestAssetAmount = histRequest.DestAssetAmount
	r.ReviewerDetails = histRequest.ReviewerDetails
	r.PreConfirmationDetails = histRequest.PreConfirmationDetails
	return
}
