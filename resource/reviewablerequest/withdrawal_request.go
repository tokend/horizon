package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateWithdrawalRequest(histRequest history.WithdrawalRequest) (
	r *regources.WithdrawalRequest, err error,
) {
	r = &regources.WithdrawalRequest{}
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
